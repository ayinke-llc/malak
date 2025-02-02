"use client";

import { useRouter } from "next/navigation";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { useMutation } from "@tanstack/react-query";
import client from "@/lib/client";
import { AxiosError } from "axios";
import { ServerAPIStatus } from "@/client/Api";
import { RiDeleteBinLine } from "@remixicon/react";

export default function DeleteDeck({ reference }: { reference: string }) {
  const router = useRouter();

  const deleteMutation = useMutation({
    mutationFn: () => client.decks.decksDelete(reference),
    onSuccess: () => {
      toast.success("Deck deleted successfully");
      router.push("/decks");
    },
    onError: (err: AxiosError<ServerAPIStatus>) => {
      toast.error(err?.response?.data?.message || "Failed to delete deck");
    },
  });

  const confirmDelete = () => {
    deleteMutation.mutate();
  };

  return (
    <AlertDialog>
      <AlertDialogTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          className="text-red-400 hover:text-red-300"
          disabled={deleteMutation.isPending}
        >
          <RiDeleteBinLine className="h-5 w-5" />
        </Button>
      </AlertDialogTrigger>

      <AlertDialogContent className="bg-white text-zinc-900 border-zinc-300 dark:bg-zinc-900 dark:text-zinc-100 dark:border-zinc-800">
        <AlertDialogHeader>
          <AlertDialogTitle>Delete deck?</AlertDialogTitle>
          <AlertDialogDescription className="text-zinc-600 dark:text-zinc-400">
            This action cannot be undone. This will permanently delete the deck
            and remove all associated data.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel
            className="bg-transparent border-zinc-300 text-zinc-900 hover:bg-zinc-200 dark:border-zinc-800 dark:text-zinc-100 dark:hover:bg-zinc-800"
            disabled={deleteMutation.isPending}
          >
            Cancel
          </AlertDialogCancel>
          <AlertDialogAction
            className="bg-red-600 text-white hover:bg-red-700"
            onClick={confirmDelete}
            disabled={deleteMutation.isPending}
          >
            {deleteMutation.isPending ? (
              <div className="flex items-center gap-2">
                <div className="h-4 w-4 animate-spin rounded-full border-2 border-zinc-900 border-t-transparent dark:border-zinc-100" />
                Deleting...
              </div>
            ) : (
              "Delete deck"
            )}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
