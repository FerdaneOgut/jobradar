FROM mcr.microsoft.com/dotnet/sdk:8.0 AS builder
WORKDIR /app
COPY TrackerService/*.csproj ./TrackerService/
RUN dotnet restore ./TrackerService/TrackerService.csproj
COPY . .
RUN dotnet publish ./TrackerService/TrackerService.csproj -c Release -o out

FROM mcr.microsoft.com/dotnet/aspnet:8.0
WORKDIR /app
COPY --from=builder /app/out .
EXPOSE 8085
CMD ["dotnet", "TrackerService.dll"]