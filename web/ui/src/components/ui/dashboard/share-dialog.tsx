import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import CreatableSelect from "@/components/ui/multi-select";
import type { OptionType } from "@/components/ui/multi-select";
import {
  RiShareBoxLine,
  RiMailLine,
  RiGlobalLine,
  RiFileCopyLine,
  RiLoader4Line,
  RiHistoryLine,
  RiSettings4Line,
  RiArrowRightLine,
} from "@remixicon/react";
import { toast } from "sonner";
import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { AccessLog } from "./access-log";
import { AccessManagement } from "./access-management";
import CopyToClipboard from 'react-copy-to-clipboard';
import { cn } from "@/lib/utils";
import { MALAK_APP_URL } from "@/lib/config";

interface ShareDialogProps {
  title: string;
  reference: string;
  token: string
}

type ShareView = "main" | "email" | "manage" | "log";

export function ShareDialog({ token, title, reference }: ShareDialogProps) {
  const [view, setView] = useState<ShareView>("main");
  const [selectedEmails, setSelectedEmails] = useState<string[]>([]);

  const fullShareLink = MALAK_APP_URL + "/shared/dashboards/" + token

  const shareDashboardMutation = useMutation({
    mutationFn: async (emails: string[]) => {
      await new Promise(resolve => setTimeout(resolve, 1000));
      return { success: true };
    },
    onSuccess: () => {
      toast.success("Dashboard shared successfully");
      setSelectedEmails([]);
      setView("main");
    },
    onError: () => {
      toast.error("Failed to share dashboard");
    }
  });

  const handleShareViaEmail = async (e: React.FormEvent) => {
    e.preventDefault();
    if (selectedEmails.length === 0) {
      toast.error("Please select at least one recipient");
      return;
    }
    await shareDashboardMutation.mutateAsync(selectedEmails);
  };

  const renderMainView = () => (
    <div className="space-y-8">
      <div className="grid grid-cols-2 gap-4">
        <button
          onClick={() => setView("email")}
          className="group relative overflow-hidden rounded-xl border bg-gradient-to-b from-muted/50 to-muted p-6 hover:shadow-md transition-all"
        >
          <div className="relative z-10 space-y-4">
            <div className="h-12 w-12 rounded-full bg-primary/10 flex items-center justify-center">
              <RiMailLine className="h-6 w-6 text-primary" />
            </div>
            <div className="space-y-2">
              <h3 className="text-xl font-semibold">Share via Email</h3>
              <p className="text-sm text-muted-foreground">
                Send an invite to specific people
              </p>
            </div>
          </div>
          <RiArrowRightLine className="absolute bottom-4 right-4 h-6 w-6 text-muted-foreground/50 transition-transform group-hover:translate-x-1" />
        </button>

        <div className="rounded-xl border bg-gradient-to-b from-muted/50 to-muted p-6">
          <div className="space-y-4">
            <div className="h-12 w-12 rounded-full bg-primary/10 flex items-center justify-center">
              <RiGlobalLine className="h-6 w-6 text-primary" />
            </div>
            <div className="space-y-2">
              <h3 className="text-xl font-semibold">Get a Link</h3>
              <p className="text-sm text-muted-foreground">
                Share with anyone who has the link
              </p>
            </div>
            <div className="w-full gap-2 mt-2">
              <CopyToClipboard
                text={fullShareLink}
                onCopy={(text, result) => {
                  if (result) {
                    toast.success('url copied. Share this with anyone');
                    return
                  }
                  toast.error('Failed to copy share url');
                }}
              >
                <Button variant="outline" className="w-full flex items-center gap-2">
                  <RiFileCopyLine className="h-5 w-5" />
                  Copy Link
                </Button>
              </CopyToClipboard>
            </div>
          </div>
        </div>
      </div >

      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-medium">Advanced Options</h3>
        </div>
        <div className="grid grid-cols-2 gap-4">
          <button
            onClick={() => setView("manage")}
            className="flex items-center gap-3 rounded-lg border p-4 hover:bg-muted/50 transition-colors"
          >
            <RiSettings4Line className="h-5 w-5 text-muted-foreground" />
            <div className="text-left">
              <h4 className="font-medium">Manage Access</h4>
              <p className="text-sm text-muted-foreground">Control who has access</p>
            </div>
          </button>
          <div className="relative flex items-center gap-3 rounded-lg border p-4 bg-gradient-to-br from-muted/30 via-muted/20 to-muted/10 overflow-hidden group cursor-not-allowed">
            <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/5 to-transparent translate-x-[-100%] group-hover:translate-x-[100%] transition-transform duration-1000" />
            <RiHistoryLine className="h-5 w-5 text-muted-foreground relative" />
            <div className="text-left relative">
              <div className="flex items-center gap-2">
                <h4 className="font-medium text-muted-foreground/80">Access Log</h4>
                <span className="text-[10px] font-medium bg-primary/10 text-primary px-2 py-0.5 rounded-full animate-pulse">Coming soon</span>
              </div>
              <p className="text-sm text-muted-foreground/70">View access history</p>
            </div>
          </div>
        </div>
      </div>
    </div >
  );

  const renderEmailView = () => (
    <div className="space-y-6">
      <form onSubmit={handleShareViaEmail} className="space-y-6">
        <div className="space-y-3">
          <label htmlFor="email-input" className="text-sm font-medium">
            Email addresses
          </label>
          <CreatableSelect
            id="email-input"
            placeholder="Type or paste email addresses"
            isMulti
            value={selectedEmails.map(email => ({ value: email, label: email } as OptionType))}
            onChange={(newValue: OptionType[]) => {
              setSelectedEmails(newValue.map(v => v.value.toString()));
            }}
            allowCustomInput={true}
            className="min-h-[42px]"
          />
          <p className="text-sm text-muted-foreground">
            Recipients will receive an email with a link to access this dashboard
          </p>
        </div>
        <div className="flex gap-3 pt-4">
          <Button
            type="button"
            variant="outline"
            className="gap-2"
            onClick={() => setView("main")}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            className="flex-1 gap-2"
            disabled={shareDashboardMutation.isPending}
          >
            {shareDashboardMutation.isPending ? (
              <>
                <RiLoader4Line className="h-5 w-5 animate-spin" />
                Sharing...
              </>
            ) : (
              <>
                <RiMailLine className="h-5 w-5" />
                Share via Email
              </>
            )}
          </Button>
        </div>
      </form>
    </div>
  );

  const viewConfig = {
    main: {
      title: "Share Dashboard",
      content: renderMainView,
      showBack: false,
    },
    email: {
      title: "Share via Email",
      content: renderEmailView,
      showBack: true,
    },
    manage: {
      title: "Manage Access",
      content: () => <AccessManagement reference={reference} />,
      showBack: true,
    },
    log: {
      title: "Access Log",
      content: () => (
        <div className="flex flex-col items-center justify-center py-12 text-center">
          <RiHistoryLine className="h-12 w-12 text-muted-foreground mb-4" />
          <h3 className="text-lg font-medium mb-2">Coming soon</h3>
          <p className="text-sm text-muted-foreground">Access log functionality will be available soon.</p>
        </div>
      ),
      showBack: true,
    },
  };

  const currentView = viewConfig[view];

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline" className="gap-2">
          <RiShareBoxLine className="h-4 w-4" />
          Share
        </Button>
      </DialogTrigger>
      <DialogContent className={cn(
        view === "email" ? "sm:max-w-md" : (view === "manage" || view === "log") ? "sm:max-w-4xl" : "sm:max-w-2xl"
      )}>
        <DialogHeader className="space-y-4">
          <div className="flex items-center gap-4">
            {currentView.showBack && (
              <button
                onClick={() => setView("main")}
                className="rounded-full p-2 hover:bg-muted transition-colors"
              >
                <RiArrowRightLine className="h-4 w-4 rotate-180" />
              </button>
            )}
            <DialogTitle className="text-2xl font-semibold">
              {currentView.title}
            </DialogTitle>
          </div>
          {view === "main" && (
            <p className="text-muted-foreground">
              Share "{title}" with others
            </p>
          )}
        </DialogHeader>

        <div className={cn(
          "mt-6",
          (view === "manage" || view === "log") && "max-h-[600px] overflow-y-auto pr-6 -mr-6"
        )}>
          {currentView.content()}
        </div>
      </DialogContent>
    </Dialog>
  );
} 
