using MongoDB.Driver;
using UserService.Models;
using UserService.DTOs;
using Microsoft.IdentityModel.Tokens;
using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using System.Text;

namespace UserService.Services;

public class UserAuthService
{
    private readonly IMongoCollection<User> _users;
    private readonly IConfiguration _config;

    public UserAuthService(IMongoClient mongoClient, IConfiguration config)
    {
        var database = mongoClient.GetDatabase("jobradar");
        _users = database.GetCollection<User>("users");
        _config = config;
    }

    public async Task<User?> GetByEmailAsync(string email)
    {
        return await _users.Find(u => u.Email == email).FirstOrDefaultAsync();
    }

    public async Task<User> CreateAsync(RegisterRequest request)
    {
        var user = new User
        {
            FullName = request.FullName,
            Email = request.Email,
            PasswordHash = BCrypt.Net.BCrypt.HashPassword(request.Password)
        };

        await _users.InsertOneAsync(user);
        return user;
    }

    public bool VerifyPassword(string password, string hash)
    {
        return BCrypt.Net.BCrypt.Verify(password, hash);
    }

    public string GenerateJwtToken(User user)
    {
        var secret = Environment.GetEnvironmentVariable("JWT_SECRET")!;
        var issuer = Environment.GetEnvironmentVariable("JWT_ISSUER")!;
        var audience = Environment.GetEnvironmentVariable("JWT_AUDIENCE")!;

        var key = new SymmetricSecurityKey(
            Encoding.UTF8.GetBytes(secret));

        var credentials = new SigningCredentials(
            key, SecurityAlgorithms.HmacSha256);

        var claims = new[]
        {
        new Claim(ClaimTypes.NameIdentifier, user.Id!),
        new Claim(ClaimTypes.Email, user.Email),
        new Claim("plan", user.Plan)
    };

        var token = new JwtSecurityToken(
            issuer: issuer,
            audience: audience,
            claims: claims,
            expires: DateTime.UtcNow.AddDays(7),
            signingCredentials: credentials
        );

        return new JwtSecurityTokenHandler().WriteToken(token);
    }
    public async Task<string> UploadCvAsync(string userId, IFormFile file)
    {
        var uploadPath = "/app/uploads";
        Directory.CreateDirectory(uploadPath);

        var extension = Path.GetExtension(file.FileName).ToLower();
        var fileName = $"{userId}{extension}";
        var filePath = Path.Combine(uploadPath, fileName);

        using (var stream = new FileStream(filePath, FileMode.Create))
        {
            await file.CopyToAsync(stream);
        }

        var update = Builders<User>.Update.Set(u => u.CvPath, filePath);
        await _users.UpdateOneAsync(u => u.Id == userId, update);

        return filePath;
    }
}