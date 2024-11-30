# 🌤️ Ecowitt to ThingSpeak

## 📄 Description

**Ecowitt to ThingSpeak** is a **Go application** that regularly (once every minute) fetches data from the [Ecowitt API](https://api.ecowitt.net/) and sends it to the [ThingSpeak](https://thingspeak.com/) platform for further analysis and visualization. The application is configured to run as a system service on Linux using `systemd` and supports cross-compilation for Windows.

![Application in Action](https://user-images.githubusercontent.com/CyberScopeToday/ecowitt_to_thingspeak/app_working.gif)

## 🚀 Requirements

- **Go**: Version 1.20 or higher
- **systemd**: For setting up the service to auto-start on Linux
- **Git**: For cloning the repository
- **Superuser Privileges**: For configuring services and installing dependencies

## 🔧 Installation and Setup

### 🧭 Step 1: Clone the Repository

First, clone the repository to your server or local machine:

```bash
git clone https://github.com/CyberScopeToday/ecowitt_to_thingspeak.git
cd ecowitt_to_thingspeak
```

### 🛠️ Step 2: Install Go

If Go is not already installed on your system, follow these steps:

#### For Linux:

1. **Download the latest version of Go:**

    ```bash
    wget https://golang.org/dl/go1.20.5.linux-amd64.tar.gz
    ```

2. **Extract the archive and install Go:**

    ```bash
    sudo tar -C /usr/local -xzf go1.20.5.linux-amd64.tar.gz
    ```

3. **Add Go to your `PATH` environment variable:**

    Open the `~/.profile` or `~/.bashrc` file and add:

    ```bash
    export PATH=$PATH:/usr/local/go/bin
    ```

4. **Apply the changes:**

    ```bash
    source ~/.profile
    ```

5. **Verify the installation:**

    ```bash
    go version
    ```

    You should see something like:

    ```
    go version go1.20.5 linux/amd64
    ```

![Go Installation](https://user-images.githubusercontent.com/CyberScopeToday/ecowitt_to_thingspeak/go_installation.png)

### 🗂️ Step 3: Create the `.env` Configuration File

Create a file named `ecowitt_to_thingspeak.env` in the root directory of the project (`~/ecowitt_to_thingspeak/`) with the following content:

```env
ECOWITT_APPLICATION_KEY=
ECOWITT_API_KEY=
ECOWITT_MAC=
THINGSPEAK_WRITE_API_KEY=
```

**Important:** Ensure that this file is **not** added to version control. Add it to `.gitignore` if using Git:

```bash
echo "ecowitt_to_thingspeak.env" >> .gitignore
```

![Environment Variables](https://user-images.githubusercontent.com/CyberScopeToday/ecowitt_to_thingspeak/env_variables.png)

### 📦 Step 4: Install Dependencies

The application uses the [`godotenv`](https://github.com/joho/godotenv) package to work with `.env` files. Install it using the following command:

```bash
go get github.com/joho/godotenv
```

### 🔨 Step 5: Build the Application

Build the executable for Linux:

```bash
go build -o ecowitt_to_thingspeak
```

For cross-compiling to Windows, run:

```bash
GOOS=windows GOARCH=amd64 go build -o ecowitt_to_thingspeak.exe
```

**Note:** Ensure that the necessary environment variables (`GOOS` and `GOARCH`) are set for the target platform.

![Build Process](https://user-images.githubusercontent.com/CyberScopeToday/ecowitt_to_thingspeak/build_process.png)

### 🖥️ Step 6: Configure systemd Service (for Linux)

#### 6.1. Create the Service File

Create a service file named `ecowitt_to_thingspeak.service` in the `/etc/systemd/system/` directory:

```bash
sudo nano /etc/systemd/system/ecowitt_to_thingspeak.service
```

Add the following content to the file:

```ini
[Unit]
Description=Ecowitt to ThingSpeak Service
After=network.target

[Service]
ExecStart=/root/ecowitt_to_thingspeak/ecowitt_to_thingspeak
WorkingDirectory=/root/ecowitt_to_thingspeak
EnvironmentFile=/root/ecowitt_to_thingspeak/ecowitt_to_thingspeak.env
Restart=always
User=root

[Install]
WantedBy=multi-user.target
```

**Important:**
- Replace `/root/ecowitt_to_thingspeak/ecowitt_to_thingspeak` with the actual path to your executable.
- Running the service as `root` is **not recommended**. It's better to create a separate user.

![Systemd Service](https://user-images.githubusercontent.com/CyberScopeToday/ecowitt_to_thingspeak/systemd_service.png)

#### 6.2. Create a User for the Service (Recommended)

1. **Create a system user:**

    ```bash
    sudo useradd -r -s /bin/false ecowitt_user
    ```

2. **Modify the service file to use the new user:**

    Open the service file:

    ```bash
    sudo nano /etc/systemd/system/ecowitt_to_thingspeak.service
    ```

    Change the `User` directive from `root` to `ecowitt_user`:

    ```ini
    User=ecowitt_user
    ```

3. **Grant the user access to the project directory:**

    ```bash
    sudo chown -R ecowitt_user:ecowitt_user /root/ecowitt_to_thingspeak
    ```

#### 6.3. Reload systemd and Start the Service

1. **Reload systemd to recognize the new service:**

    ```bash
    sudo systemctl daemon-reload
    ```

2. **Start the service:**

    ```bash
    sudo systemctl start ecowitt_to_thingspeak.service
    ```

3. **Enable the service to start on boot:**

    ```bash
    sudo systemctl enable ecowitt_to_thingspeak.service
    ```

4. **Check the service status:**

    ```bash
    sudo systemctl status ecowitt_to_thingspeak.service
    ```

    You should see something like:

    ```
    ● ecowitt_to_thingspeak.service - Ecowitt to ThingSpeak Service
       Loaded: loaded (/etc/systemd/system/ecowitt_to_thingspeak.service; enabled; vendor preset: enabled)
       Active: active (running) since Mon 2024-04-27 12:00:01 UTC; 5s ago
     Main PID: 12345 (ecowitt_to_t)
        Tasks: 1 (limit: 4915)
       Memory: 10.0M
          CPU: 0.05s
       CGroup: /system.slice/ecowitt_to_thingspeak.service
               └─12345 /root/ecowitt_to_thingspeak/ecowitt_to_thingspeak
    ```

    ![Service Status](https://user-images.githubusercontent.com/CyberScopeToday/ecowitt_to_thingspeak/service_status.png)

#### 6.4. View Service Logs

To view the service logs, use:

```bash
journalctl -u ecowitt_to_thingspeak.service -f
```

![Service Logs](https://user-images.githubusercontent.com/CyberScopeToday/ecowitt_to_thingspeak/service_logs.png)

### 🖱️ Step 7: Manually Run the Application (Optional)

If you prefer to run the application manually without using `systemd`, execute:

```bash
./ecowitt_to_thingspeak
```

## 💻 Cross-Compilation for Windows

To create an executable for Windows, follow these steps:

1. **Navigate to the project directory:**

    ```bash
    cd ~/ecowitt_to_thingspeak
    ```

2. **Build the executable for Windows (64-bit):**

    ```bash
    GOOS=windows GOARCH=amd64 go build -o ecowitt_to_thingspeak.exe
    ```

3. **Copy `ecowitt_to_thingspeak.exe` and the `ecowitt_to_thingspeak.env` file to your Windows machine.**

4. **Run the application on Windows:**

    Open Command Prompt or PowerShell, navigate to the directory containing the executable, and run:

    ```cmd
    ecowitt_to_thingspeak.exe
    ```

    Ensure that the application runs correctly and sends data to ThingSpeak.

![Windows Application](https://user-images.githubusercontent.com/CyberScopeToday/ecowitt_to_thingspeak/windows_application.png)

### Setting Up Auto-Start on Windows

To make the application start automatically when Windows boots:

1. **Create a Shortcut for the Application:**
    - Navigate to the directory containing `ecowitt_to_thingspeak.exe`.
    - Right-click on `ecowitt_to_thingspeak.exe` and select **Create shortcut**.

2. **Move the Shortcut to the Startup Folder:**
    - Press `Win + R`, type `shell:startup`, and press `Enter`.
    - Move the created shortcut into the opened Startup folder.

![Windows Startup](https://user-images.githubusercontent.com/CyberScopeToday/ecowitt_to_thingspeak/windows_startup.png)

## 🔒 Security

### Restrict Access to the `.env` File

Ensure that the `.env` file is only readable by the owner:

```bash
chmod 600 /root/ecowitt_to_thingspeak/ecowitt_to_thingspeak.env
```

### Security Recommendations

- **Do not run the service as `root` unless necessary.** Create a dedicated user with minimal permissions.
- **Do not share your `.env` file publicly.** Ensure it is added to `.gitignore`.
- **Regularly update dependencies** and check for vulnerabilities.

## 📝 Logging

The application uses Go's standard `log` package to output information. Logs can be viewed through `journalctl` on Linux or redirected to a file when running manually.

![Logging](https://user-images.githubusercontent.com/CyberScopeToday/ecowitt_to_thingspeak/logging.png)

## 📈 Monitoring and Notifications

It is recommended to add a notification system (e.g., via Telegram or Email) to alert you of critical errors or service failures.

![Monitoring](https://user-images.githubusercontent.com/CyberScopeToday/ecowitt_to_thingspeak/monitoring.png)

## 🔄 Updating the Application

To update the application, follow these steps:

1. **Pull the latest changes from the repository:**

    ```bash
    git pull origin main
    ```

2. **Install any new dependencies:**

    ```bash
    go mod tidy
    ```

3. **Rebuild the executable:**

    ```bash
    go build -o ecowitt_to_thingspeak
    ```

4. **Restart the service:**

    ```bash
    sudo systemctl restart ecowitt_to_thingspeak.service
    ```

![Update Process](https://user-images.githubusercontent.com/CyberScopeToday/ecowitt_to_thingspeak/update_process.png)

## 🤝 Contributing and Support

If you have suggestions for improvements or find a bug, please create an [Issue](https://github.com/CyberScopeToday/ecowitt_to_thingspeak/issues) in the repository.

![Contribute](https://user-images.githubusercontent.com/CyberScopeToday/ecowitt_to_thingspeak/contribute.png)

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---
