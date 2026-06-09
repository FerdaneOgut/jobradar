using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;

namespace TrackerService.Models;

public class Application
{
    [BsonId]
    [BsonRepresentation(BsonType.ObjectId)]
    public string? Id { get; set; }

    [BsonElement("userId")]
    public string UserId { get; set; } = string.Empty;

    [BsonElement("jobId")]
    public string JobId { get; set; } = string.Empty;

    [BsonElement("jobTitle")]
    public string JobTitle { get; set; } = string.Empty;

    [BsonElement("company")]
    public string Company { get; set; } = string.Empty;

    [BsonElement("jobUrl")]
    public string JobUrl { get; set; } = string.Empty;

    [BsonElement("status")]
    public string Status { get; set; } = "saved";

    [BsonElement("notes")]
    public string Notes { get; set; } = string.Empty;

    [BsonElement("appliedAt")]
    public DateTime? AppliedAt { get; set; }

    [BsonElement("followUpSentAt")]
    public DateTime? FollowUpSentAt { get; set; }

    [BsonElement("createdAt")]
    public DateTime CreatedAt { get; set; } = DateTime.UtcNow;

    [BsonElement("updatedAt")]
    public DateTime UpdatedAt { get; set; } = DateTime.UtcNow;
}