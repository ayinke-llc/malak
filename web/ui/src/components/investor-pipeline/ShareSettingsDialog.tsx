import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { RiFileCopyLine, RiRefreshLine } from "@remixicon/react";
import { toast } from "sonner";

interface ShareSettings {
  isEnabled: boolean;
  shareLink: string;
  requireEmail: boolean;
  requirePassword: boolean;
  password?: string;
}

interface ShareSettingsDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  settings: ShareSettings;
  onSettingsChange: (settings: ShareSettings) => void;
}

// Mock function to generate share link
const generateShareLink = () => {
  return `https://your-domain.com/share/${Math.random().toString(36).substring(2, 15)}`;
};

function EmailSection({ settings, onSettingsChange }: { 
  settings: ShareSettings; 
  onSettingsChange: (settings: ShareSettings) => void;
}) {
  return (
    <div className="flex items-center justify-between">
      <div className="space-y-0.5">
        <Label>Require Email</Label>
        <p className="text-sm text-muted-foreground">
          Viewers must verify their email before accessing
        </p>
      </div>
      <Switch
        checked={settings.requireEmail}
        onCheckedChange={(checked) =>
          onSettingsChange({ ...settings, requireEmail: checked })
        }
      />
    </div>
  );
}

function PasswordSection({ settings, onSettingsChange }: {
  settings: ShareSettings;
  onSettingsChange: (settings: ShareSettings) => void;
}) {
  return (
    <>
      <div className="flex items-center justify-between">
        <div className="space-y-0.5">
          <Label>Password Protection</Label>
          <p className="text-sm text-muted-foreground">
            Require a password to access the pipeline
          </p>
        </div>
        <Switch
          checked={settings.requirePassword}
          onCheckedChange={(checked) =>
            onSettingsChange({ ...settings, requirePassword: checked })
          }
        />
      </div>

      {settings.requirePassword && (
        <div className="space-y-2 pt-2">
          <Label>Set Password</Label>
          <Input
            type="password"
            placeholder="Enter password for protection"
            value={settings.password || ""}
            onChange={(e) =>
              onSettingsChange({ ...settings, password: e.target.value })
            }
          />
          <p className="text-sm text-muted-foreground">
            Make sure to share this password securely with your viewers
          </p>
        </div>
      )}
    </>
  );
}

function ShareLinkSection({ settings, onSettingsChange }: {
  settings: ShareSettings;
  onSettingsChange: (settings: ShareSettings) => void;
}) {
  const handleCopyLink = () => {
    navigator.clipboard.writeText(settings.shareLink);
    toast.success("Link copied to clipboard");
  };

  const handleResetLink = () => {
    onSettingsChange({
      ...settings,
      shareLink: generateShareLink(),
    });
    toast.success("New share link generated");
  };

  return (
    <div className="space-y-2">
      <Label>Share Link</Label>
      <div className="flex gap-2">
        <Input
          readOnly
          value={settings.shareLink}
          className="flex-1"
        />
        <Button
          variant="outline"
          size="icon"
          onClick={handleCopyLink}
          title="Copy link"
        >
          <RiFileCopyLine className="h-4 w-4" />
        </Button>
        <Button
          variant="outline"
          size="icon"
          onClick={handleResetLink}
          title="Reset link"
        >
          <RiRefreshLine className="h-4 w-4" />
        </Button>
      </div>
      <p className="text-sm text-muted-foreground">
        {!settings.requireEmail && !settings.requirePassword
          ? "Anyone with this link can view your pipeline"
          : `Viewers will need to ${[
              settings.requireEmail && "verify their email",
              settings.requirePassword && "enter the password"
            ].filter(Boolean).join(" and ")} to access`}
      </p>
    </div>
  );
}

export function ShareSettingsDialog({
  open,
  onOpenChange,
  settings,
  onSettingsChange,
}: ShareSettingsDialogProps) {
  const handleToggleShare = (enabled: boolean) => {
    onSettingsChange({
      ...settings,
      isEnabled: enabled,
      shareLink: enabled ? generateShareLink() : "",
    });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Share Settings</DialogTitle>
          <DialogDescription>
            Configure how you want to share your investor pipeline
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6 py-4">
          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <Label>Enable Sharing</Label>
              <p className="text-sm text-muted-foreground">
                Allow others to view your investor pipeline via link
              </p>
            </div>
            <Switch
              checked={settings.isEnabled}
              onCheckedChange={handleToggleShare}
            />
          </div>

          {settings.isEnabled && (
            <div className="space-y-4">
              <EmailSection settings={settings} onSettingsChange={onSettingsChange} />
              <PasswordSection settings={settings} onSettingsChange={onSettingsChange} />
              <ShareLinkSection settings={settings} onSettingsChange={onSettingsChange} />
            </div>
          )}
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Close
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
} 