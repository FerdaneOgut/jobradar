FROM mcr.microsoft.com/dotnet/sdk:8.0 AS builder
WORKDIR /app
COPY UserService/*.csproj ./UserService/
RUN dotnet restore ./UserService/UserService.csproj
COPY . .
RUN dotnet publish ./UserService/UserService.csproj -c Release -o out

FROM mcr.microsoft.com/dotnet/aspnet:8.0
WORKDIR /app
COPY --from=builder /app/out .
EXPOSE 8084
CMD ["dotnet", "UserService.dll"]