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
import {
  useReactTable,
  getCoreRowModel,
  flexRender,
  createColumnHelper,
} from '@tanstack/react-table';
import { useVirtualizer } from '@tanstack/react-virtual';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import CopyToClipboard from 'react-copy-to-clipboard';
import { useMutation } from '@tanstack/react-query';
import { IMPORT_CONTACTS_MUTATION } from '@/lib/query-constants';
import client from '@/lib/client';

const MAX_FILE_SIZE = 500 * 1024; // 500KB in bytes

interface FileWithSize extends File {
  size: number;
  type: string;
}

const schema = yup.object().shape({
  file: yup
    .mixed<FileWithSize>()
    .required("File is required")
    .test("fileSize", "File size must be less than 500KB", (value): value is FileWithSize => {
      if (!value) return false;
      return (value as FileWithSize).size <= MAX_FILE_SIZE;
    })
    .test("fileType", "Only CSV files are allowed", (value): value is FileWithSize => {
      if (!value) return false;
      return (value as FileWithSize).type === "text/csv";
    }),
});

type FormData = yup.InferType<typeof schema>;

const truncateText = (text: string, maxLength: number) => {
  if (!text) return '-';
  if (text.length <= maxLength) return text;
  return text.slice(0, maxLength) + "...";
};

type Contact = {
  email: string;
  first_name: string;
  last_name: string;
  company: string;
  notes: string;
};

export default function CSVUploadModal() {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [uploadSuccess, setUploadSuccess] = useState(false);
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [headers, setHeaders] = useState<string[]>([]);
  const [tableContainerRef, setTableContainerRef] = useState<HTMLDivElement | null>(null);
  const [isOpen, setIsOpen] = useState(false);

  const { mutate: importContacts, isPending } = useMutation({
    mutationKey: [IMPORT_CONTACTS_MUTATION],
    mutationFn: async (contacts: Contact[]) => {
      const response = await client.contacts.batchCreate({ contacts });
      return response.data;
    },
    onSuccess: () => {
      toast.success('Contacts imported successfully');
      setIsOpen(false);
      setContacts([]);
      setSelectedFile(null);
      setUploadSuccess(false);
    },
    onError: (error: Error) => {
      toast.error(`Failed to import contacts: ${error.message}`);
    },
  });

  const columnHelper = createColumnHelper<Contact>();

  const columns = headers.map(header => 
    columnHelper.accessor(header as keyof Contact, {
      header: () => header.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase()),
      cell: info => {
        const value = info.getValue();
        if (header === 'notes') {
          return (
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <div className="cursor-pointer max-w-[150px]">
                    {truncateText(value || '', 50)}
                  </div>
                </TooltipTrigger>
                <TooltipContent side="bottom" className="max-w-sm p-4">
                  <p className="whitespace-pre-wrap break-words">{value || '-'}</p>
                  <CopyToClipboard 
                    text={value || ''}
                    onCopy={() => toast.success('Copied to clipboard')}
                  >
                    <Button
                      variant="ghost"
                      size="sm"
                      className="mt-2 h-8 w-full text-xs"
                    >
                      Copy to clipboard
                    </Button>
                  </CopyToClipboard>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          );
        }
        return value || '-';
      },
      size: header === 'notes' ? 150 : 150,
    })
  );

  const table = useReactTable({
    data: contacts,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  const rowVirtualizer = useVirtualizer({
    count: contacts.length,
    getScrollElement: () => tableContainerRef,
    estimateSize: () => 40,
    overscan: 5,
  });

  const virtualRows = rowVirtualizer.getVirtualItems();

  const {
    register,
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
      const rows = text.split('\n').filter(row => row.trim());
      
      if (rows.length < 2) {
        toast.error('CSV file must contain headers and at least one contact');
        return;
      }

      if (rows.length > 101) { // 101 because first row is header
        toast.error('CSV file cannot contain more than 100 contacts');
        return;
      }

      const csvHeaders = rows[0].split(',').map(header => header.trim().toLowerCase());
      setHeaders(csvHeaders);
      
      const requiredHeaders = ['email', 'first_name', 'last_name', 'company', 'notes'];
      const missingHeaders = requiredHeaders.filter(header => !csvHeaders.includes(header));
      
      if (missingHeaders.length > 0) {
        toast.error(`Missing required headers: ${missingHeaders.join(', ')}`);
        return;
      }

      const processedContacts = rows.slice(1)
        .filter(row => row.trim())
        .map((row, index) => {
          const values = row.split(',').map(value => value.trim());
          const contact: Record<string, string> = {};
          
          csvHeaders.forEach((header, colIndex) => {
            contact[header] = values[colIndex] || '';
          });

          return contact as Contact;
        });

      const invalidEmails = processedContacts.filter(contact => !contact.email);
      if (invalidEmails.length > 0) {
        toast.error(`Found ${invalidEmails.length} contacts missing required email field`);
        return;
      }

      setContacts(processedContacts);
      toast.success(`Successfully processed ${processedContacts.length} contacts`);
    };

    reader.readAsText(file);
  };

  const handleImportSubmit = () => {
    if (contacts.length === 0) {
      toast.error('No contacts to import');
      return;
    }
    importContacts(contacts);
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setValue('file', file);
      setSelectedFile(file);
      setUploadSuccess(false);
      setContacts([]);
      setHeaders([]);
      
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
      setContacts([]);
      setHeaders([]);
      
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

  const paddingTop = virtualRows.length > 0 ? virtualRows?.[0]?.start || 0 : 0;
  const paddingBottom =
    virtualRows.length > 0
      ? rowVirtualizer.getTotalSize() - (virtualRows?.[virtualRows.length - 1]?.end || 0)
      : 0;

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
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
      <DialogContent className="sm:max-w-4xl">
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

          {contacts.length > 0 && (
            <div className="mt-6">
              <h4 className="text-sm font-medium mb-2">Preview ({contacts.length} contacts)</h4>
              <div 
                ref={setTableContainerRef}
                className="virtual-table-container border rounded-md" 
                style={{ height: '400px', overflow: 'auto' }}
              >
                <Table>
                  <TableHeader className="sticky top-0 bg-background z-10">
                    {table.getHeaderGroups().map(headerGroup => (
                      <TableRow key={headerGroup.id}>
                        {headerGroup.headers.map(header => (
                          <TableHead 
                            key={header.id}
                            style={{ width: header.getSize() }}
                            className="whitespace-nowrap px-4 py-2"
                          >
                            {flexRender(
                              header.column.columnDef.header,
                              header.getContext()
                            )}
                          </TableHead>
                        ))}
                      </TableRow>
                    ))}
                  </TableHeader>
                  <TableBody>
                    {paddingTop > 0 && (
                      <tr>
                        <td style={{ height: `${paddingTop}px` }} />
                      </tr>
                    )}
                    {virtualRows.map((virtualRow) => {
                      const row = table.getRowModel().rows[virtualRow.index];
                      return (
                        <TableRow
                          key={row.id}
                          className="hover:bg-hover"
                        >
                          {row.getVisibleCells().map(cell => (
                            <TableCell 
                              key={cell.id}
                              style={{ width: cell.column.getSize() }}
                              className="px-4 py-2"
                            >
                              {flexRender(cell.column.columnDef.cell, cell.getContext())}
                            </TableCell>
                          ))}
                        </TableRow>
                      );
                    })}
                    {paddingBottom > 0 && (
                      <tr>
                        <td style={{ height: `${paddingBottom}px` }} />
                      </tr>
                    )}
                  </TableBody>
                </Table>
              </div>
              <div className="mt-4 flex justify-end">
                <Button
                  type="button"
                  onClick={handleImportSubmit}
                  className="gap-2"
                  disabled={isPending}
                >
                  {isPending ? 'Importing...' : `Import ${contacts.length} Contacts`}
                </Button>
              </div>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
} 