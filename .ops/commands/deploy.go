package commands

import (
	"context"
	"fmt"
	"os"

	"lesiw.io/command"
	"lesiw.io/command/sys"
)

func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

const (
	piHost      = "kevin@jor"
	piDir       = "/home/kevin/Downloads/crons"
	serviceName = "crons"
)

var serviceFile = fmt.Sprintf(`[Unit]
Description=Crons Light Scheduler
After=network.target

[Service]
ExecStart=%s/crons
WorkingDirectory=%s
EnvironmentFile=%s/.env
Restart=always

[Install]
WantedBy=multi-user.target
`, piDir, piDir, piDir)

func (Ops) Deploy() error {
	ctx := context.Background()
	sh := command.Shell(sys.Machine(), "go", "scp", "ssh")

	fmt.Println("==> Building for Raspberry Pi (linux/arm64)...")
	var ops Ops
	if err := ops.Build(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	fmt.Println("==> Creating directory on Pi...")
	if err := sh.Exec(ctx, "ssh", piHost, "mkdir", "-p", piDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Println("==> Stopping existing service...")
	// Ignore error - service may not exist on first deploy
	_ = sh.Exec(ctx, "ssh", piHost, "sudo", "systemctl", "stop", serviceName)

	fmt.Println("==> Copying binary to Pi...")
	if err := sh.Exec(ctx, "scp", "crons", piHost+":"+piDir+"/crons"); err != nil {
		return fmt.Errorf("failed to copy binary: %w", err)
	}
	if err := sh.Exec(ctx, "ssh", piHost, "chmod", "+x", piDir+"/crons"); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}

	fmt.Println("==> Copying .env to Pi...")
	if err := sh.Exec(ctx, "scp", ".env", piHost+":"+piDir+"/.env"); err != nil {
		return fmt.Errorf("failed to copy .env: %w", err)
	}

	fmt.Println("==> Setting up systemd service...")
	// Write service file locally then copy
	serviceFilePath := serviceName + ".service"
	if err := writeFile(serviceFilePath, serviceFile); err != nil {
		return fmt.Errorf("failed to write service file locally: %w", err)
	}
	if err := sh.Exec(ctx, "scp", serviceFilePath, piHost+":/tmp/"+serviceFilePath); err != nil {
		return fmt.Errorf("failed to copy service file: %w", err)
	}
	if err := sh.Exec(ctx, "ssh", piHost, "sudo", "mv", "/tmp/"+serviceFilePath, "/etc/systemd/system/"+serviceFilePath); err != nil {
		return fmt.Errorf("failed to move service file: %w", err)
	}

	fmt.Println("==> Enabling and restarting service...")
	if err := sh.Exec(ctx, "ssh", piHost, "sudo", "systemctl", "daemon-reload"); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}
	if err := sh.Exec(ctx, "ssh", piHost, "sudo", "systemctl", "enable", serviceName); err != nil {
		return fmt.Errorf("failed to enable service: %w", err)
	}
	if err := sh.Exec(ctx, "ssh", piHost, "sudo", "systemctl", "restart", serviceName); err != nil {
		return fmt.Errorf("failed to restart service: %w", err)
	}

	fmt.Println("==> Checking service status...")
	if err := sh.Exec(ctx, "ssh", piHost, "sudo", "systemctl", "status", serviceName, "--no-pager"); err != nil {
		return fmt.Errorf("service not running: %w", err)
	}

	fmt.Println("==> Deploy complete!")
	return nil
}
