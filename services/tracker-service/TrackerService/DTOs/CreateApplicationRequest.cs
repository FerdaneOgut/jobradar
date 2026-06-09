namespace TrackerService.DTOs;

public class CreateApplicationRequest
{
    public string JobId { get; set; } = string.Empty;
    public string JobTitle { get; set; } = string.Empty;
    public string Company { get; set; } = string.Empty;
    public string JobUrl { get; set; } = string.Empty;
    public string Notes { get; set; } = string.Empty;
}