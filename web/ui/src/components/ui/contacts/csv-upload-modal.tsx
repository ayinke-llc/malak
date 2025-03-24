"use client"

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { RiAddLine, RiUploadCloud2Line } from "@remixicon/react";
import { useState } from "react";
import { toast } from "sonner";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";

const MAX_FILE_SIZE = 500 * 1024; // 500KB in bytes

const schema = yup.object().shape({
  file: yup
    .mixed()
    .required("File is required")
    .test("fileSize", "File size must be less than 500KB", (value) => {
      if (!value) return false;
      return value.size <= MAX_FILE_SIZE;
    })
    .test("fileType", "Only CSV files are allowed", (value) => {
      if (!value) return false;
      return value.type === "text/csv";
    }),
});

type FormData = yup.InferType<typeof schema>;

const truncateText = (text: string, maxLength: number) => {
  if (text.length <= maxLength) return text;
  return text.slice(0, maxLength) + "...";
};

export default function CSVUploadModal() {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [uploadSuccess, setUploadSuccess] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
    reset,
  } = useForm<FormData>({
    resolver: yupResolver(schema),
  });

  const processCSV = (file: File) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      const text = e.target?.result as string;
      // Split by newlines and filter out empty lines
      const rows = text.split('\n').filter(row => row.trim());
      
      if (rows.length < 2) {
        toast.error('CSV file must contain headers and at least one contact');
        return;
      }

      const headers = rows[0].split(',').map(header => header.trim().toLowerCase());
      
      // Validate headers
      const requiredHeaders = ['email', 'first_name', 'last_name', 'company', 'notes'];
      const missingHeaders = requiredHeaders.filter(header => !headers.includes(header));
      
      if (missingHeaders.length > 0) {
        toast.error(`Missing required headers: ${missingHeaders.join(', ')}`);
        return;
      }

      // Process rows (skip header row)
      const contacts = rows.slice(1)
        .filter(row => row.trim()) // Filter out empty rows
        .map((row, index) => {
          const values = row.split(',').map(value => value.trim());
          const contact: Record<string, string> = {};
          
          headers.forEach((header, colIndex) => {
            // Handle missing values in the row
            contact[header] = values[colIndex] || '';
          });

          return contact;
        });

      // Validate required email field
      const invalidEmails = contacts.filter(contact => !contact.email);
      if (invalidEmails.length > 0) {
        toast.error(`Found ${invalidEmails.length} contacts missing required email field`);
        return;
      }

      // Here you would typically send the contacts to your API
      console.log('Processed contacts:', contacts);
      toast.success(`Successfully processed ${contacts.length} contacts`);
    };

    reader.readAsText(file);
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setValue('file', file);
      setSelectedFile(file);
      setUploadSuccess(false);
      
      try {
        await schema.validateAt('file', { file });
        processCSV(file);
        setUploadSuccess(true);
      } catch (error) {
        if (error instanceof yup.ValidationError) {
          toast.error(error.message);
        }
      }
    }
  };

  const handleDrop = async (e: React.DragEvent<HTMLLabelElement>) => {
    e.preventDefault();
    const file = e.dataTransfer.files?.[0];
    if (file) {
      setValue('file', file);
      setSelectedFile(file);
      setUploadSuccess(false);
      
      try {
        await schema.validateAt('file', { file });
        processCSV(file);
        setUploadSuccess(true);
      } catch (error) {
        if (error instanceof yup.ValidationError) {
          toast.error(error.message);
        }
      }
    }
  };

  const handleDragOver = (e: React.DragEvent<HTMLLabelElement>) => {
    e.preventDefault();
  };

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button
          type="button"
          variant="default"
          className="whitespace-nowrap gap-1"
        >
          <RiAddLine />
          Import CSV
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="text-foreground">Import Contacts from CSV</DialogTitle>
          <DialogDescription className="text-muted-foreground">
            Upload a CSV file containing your contacts. The CSV should have the following columns:
            email (required), first_name, last_name, company, notes
          </DialogDescription>
        </DialogHeader>

        <div className="mt-6">
          <Label htmlFor="csv" className="text-foreground">CSV File</Label>
          <div className="flex items-center justify-center w-full">
            <label
              htmlFor="csv"
              className="flex flex-col items-center justify-center w-full h-32 border-2 border-dashed rounded-lg cursor-pointer border-input bg-background/50 hover:bg-hover transition-colors duration-200"
              onDrop={handleDrop}
              onDragOver={handleDragOver}
            >
              <div className="flex flex-col items-center justify-center pt-5 pb-6">
                <RiUploadCloud2Line className="w-8 h-8 mb-3 text-muted-foreground" />
                {selectedFile ? (
                  <div className="text-center">
                    <TooltipProvider>
                      <Tooltip>
                        <TooltipTrigger asChild>
                          <p className="mb-2 text-sm text-muted-foreground">
                            Selected: <span className="font-medium text-foreground">{truncateText(selectedFile.name, 40)}</span>
                          </p>
                        </TooltipTrigger>
                        <TooltipContent>
                          <p>{selectedFile.name}</p>
                        </TooltipContent>
                      </Tooltip>
                    </TooltipProvider>
                    <p className="text-xs text-muted-foreground">Click or drag to change file</p>
                  </div>
                ) : (
                  <>
                    <p className="mb-2 text-sm text-muted-foreground">
                      <span className="font-semibold">Click to upload</span> or drag and drop
                    </p>
                    <p className="text-xs text-muted-foreground">CSV files only (MAX. 500KB)</p>
                  </>
                )}
              </div>
              <input
                id="csv"
                type="file"
                className="hidden"
                accept=".csv"
                {...register('file')}
                onChange={handleFileChange}
              />
            </label>
            {errors.file && (
              <p className="mt-1 text-xs text-red-600 dark:text-red-500">
                <span className="font-medium">{errors.file.message}</span>
              </p>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
} 