using MongoDB.Driver;
using TrackerService.Models;
using TrackerService.DTOs;

namespace TrackerService.Services;

public class ApplicationService
{
    private readonly IMongoCollection<Application> _applications;
    private readonly ILogger<ApplicationService> _logger;

    public ApplicationService(IMongoClient mongoClient, ILogger<ApplicationService> logger)
    {
        var database = mongoClient.GetDatabase("jobradar");
        _applications = database.GetCollection<Application>("applications");
        _logger = logger;
    }

    public async Task<Application> CreateAsync(string userId, CreateApplicationRequest request)
    {
        var application = new Application
        {
            UserId = userId,
            JobId = request.JobId,
            JobTitle = request.JobTitle,
            Company = request.Company,
            JobUrl = request.JobUrl,
            Notes = request.Notes,
            Status = "saved"
        };

        await _applications.InsertOneAsync(application);
        return application;
    }

    public async Task<List<Application>> GetByUserIdAsync(string userId)
    {
        return await _applications
            .Find(a => a.UserId == userId)
            .SortByDescending(a => a.CreatedAt)
            .ToListAsync();
    }

    public async Task<Application?> GetByIdAsync(string id)
    {
        return await _applications.Find(a => a.Id == id).FirstOrDefaultAsync();
    }

    public async Task<Application?> UpdateAsync(string id, string userId, UpdateApplicationRequest request)
    {
        var application = await GetByIdAsync(id);
        if (application == null || application.UserId != userId)
            return null;

        var update = Builders<Application>.Update
            .Set(a => a.UpdatedAt, DateTime.UtcNow);

        if (request.Status != null)
        {
            update = update.Set(a => a.Status, request.Status);

            if (request.Status == "applied")
                update = update.Set(a => a.AppliedAt, DateTime.UtcNow);
        }

        if (request.Notes != null)
            update = update.Set(a => a.Notes, request.Notes);

        await _applications.UpdateOneAsync(a => a.Id == id, update);
        return await GetByIdAsync(id);
    }

    public async Task<bool> DeleteAsync(string id, string userId)
    {
        var result = await _applications.DeleteOneAsync(a => a.Id == id && a.UserId == userId);
        return result.DeletedCount > 0;
    }

    public async Task<List<Application>> GetPendingFollowUpsAsync()
    {
        var sevenDaysAgo = DateTime.UtcNow.AddDays(-7);
        return await _applications
            .Find(a => a.Status == "applied"
                && a.AppliedAt <= sevenDaysAgo
                && a.FollowUpSentAt == null)
            .ToListAsync();
    }
}