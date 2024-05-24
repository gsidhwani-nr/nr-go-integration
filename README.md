# nr-go-integration

# Service Demonstration with New Relic

This repository contains three Go applications: `Service1`, `Service2`, and `ClientApp`. The `ClientApp` calls `Service1`, which in turn calls `Service2`. The JSON response from `Service2` is processed by `Service1` through a goroutine to extract book titles, and finally, the titles are returned to the client. New Relic is used to track the services, goroutines, and to collect custom attributes for monitoring and tracing.

## Prerequisites

- Go 1.16+
- New Relic Account and License Key
- `Gin` Web Framework

## Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/your-repo/nr-go-integration.git
   cd nr-go-integration
   ```
2. **Set the New Relic License Key**

Make sure to set the New Relic License Key as an environment variable before running the services:   
   ```bash
  export NEW_RELIC_LICENSE_KEY="your_new_relic_license_key"
   ```
3. **Building and Running the Services. Please open a separate bash terminal**
   - Service2
     - Service2 provides a list of books in JSON format.
     ```bash
     go build -o service2 ./service2.go
     ./service2
     ```

     - Service1
     - Service1 calls Service2, processes the JSON response to extract book titles using a goroutine, and returns the titles to the client.
     ```bash
     go build -o service1 ./service1.go
     ./service2
     ```

     - ClientApp
     - ClientApp makes multiple calls to Service1 and logs the book titles returned by Service1.
     ```bash
     go build -o clientapp ./clientApp.go
     ./clientapp
     ```
## Summary
- ClientApp makes calls to Service1, which in turn calls Service2.
- Service1 processes the JSON response from Service2 in a goroutine and returns the book titles.
- New Relic is used to monitor and trace the services and goroutines, and to collect custom attributes for monitoring.

- <img width="880" alt="image" src="https://github.com/gsidhwani-nr/nr-go-integration/assets/113113837/cab3f68b-28be-467e-997b-3ee3d11f8983">
<img width="1680" alt="image" src="https://github.com/gsidhwani-nr/nr-go-integration/assets/113113837/1d72bda8-fe6f-4bce-bf42-54482efdc68c">
<img width="1680" alt="image" src="https://github.com/gsidhwani-nr/nr-go-integration/assets/113113837/3f31858f-3247-476c-9720-93dc11a55993">

## Troubleshooting 

if you face some project initiliazation error, please re-initilaize the project.

```bash
go mod init gsidhwani-nr/nr-go-integration 
go mod tidy
```
