import { useState, useRef, useEffect } from "react";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  RiMailLine,
  RiPhoneLine,
  RiCalendarLine,
  RiFileTextLine,
  RiMoneyDollarCircleLine,
  RiTeamLine,
  RiCloseLine,
  RiTimeLine,
  RiAddLine,
  RiArrowRightLine,
  RiEditLine,
  RiDeleteBinLine,
  RiDownloadLine,
  RiUploadLine,
  RiFile3Line,
  RiFileExcelLine,
  RiImageLine,
  RiMoreLine,
  RiUploadCloud2Line,
  RiStarFill,
  RiStarLine,
  RiArchiveFill,
} from "@remixicon/react";
import { Skeleton } from "@/components/ui/skeleton";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { toast } from "sonner";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { NumericFormat } from "react-number-format";

interface Contact {
  name: string;
  image: string;
}

interface Card {
  id: string;
  title: string;
  amount: string;
  stage: string;
  dueDate: string;
  contact: Contact;
  checkSize: string;
  initialContactDate: string;
  isLeadInvestor: boolean;
  rating: number;
}

interface Note {
  id: string;
  title: string;
  content: string;
  createdAt: string;
  updatedAt: string;
}

interface Activity {
  id: string;
  type: 'email' | 'meeting' | 'note';
  title: string;
  description: string;
  timestamp: string;
  content?: string;
}

interface Document {
  id: string;
  name: string;
  type: 'pdf' | 'excel' | 'image' | 'other';
  size: number;
  uploadedAt: Date;
  uploadedBy: string;
}

interface InvestorDetailsDrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  investor: Card | null;
  isArchived?: boolean;
}

const TOTAL_ACTIVITIES_LIMIT = 250;
const ACTIVITIES_PER_PAGE = 25;

// Mock function to fetch activities
const fetchActivities = async (page: number, investorId: string): Promise<Activity[]> => {
  // Simulate API delay
  await new Promise(resolve => setTimeout(resolve, 1000));
  
  const startIndex = page * ACTIVITIES_PER_PAGE;
  // If we've reached the limit, return empty array
  if (startIndex >= TOTAL_ACTIVITIES_LIMIT) {
    return [];
  }
  
  // Calculate how many items to generate (handle last page)
  const itemsToGenerate = Math.min(
    ACTIVITIES_PER_PAGE,
    TOTAL_ACTIVITIES_LIMIT - startIndex
  );
  
  // Generate activities
  return Array.from({ length: itemsToGenerate }, (_, i) => ({
    id: `${page}-${i}-${Math.random()}`,
    type: ['email', 'meeting', 'note'][Math.floor(Math.random() * 3)] as Activity['type'],
    title: `Activity ${startIndex + i + 1}`,
    description: `Description for activity ${startIndex + i + 1}`,
    timestamp: new Date(Date.now() - (startIndex + i) * 24 * 60 * 60 * 1000).toISOString(),
    content: Math.random() > 0.5 ? `Content for activity ${startIndex + i + 1}` : undefined,
  }));
};

// Mock documents data
const mockDocuments: Document[] = [
  {
    id: '1',
    name: 'Financial Report Q4 2023.pdf',
    type: 'pdf',
    size: 2500000,
    uploadedAt: new Date('2024-01-15'),
    uploadedBy: 'John Doe'
  },
  {
    id: '2',
    name: 'Investment Metrics.xlsx',
    type: 'excel',
    size: 1800000,
    uploadedAt: new Date('2024-01-20'),
    uploadedBy: 'Jane Smith'
  },
  {
    id: '3',
    name: 'Company Logo.png',
    type: 'image',
    size: 500000,
    uploadedAt: new Date('2024-01-25'),
    uploadedBy: 'John Doe'
  }
];

function ActivitySkeleton() {
  return (
    <div className="relative">
      <div className="absolute -left-[27px] bg-background p-1 border rounded-full">
        <Skeleton className="h-4 w-4 rounded-full" />
      </div>
      <div className="bg-card rounded-lg p-4 border">
        <div className="flex items-center justify-between mb-2">
          <Skeleton className="h-4 w-32" />
          <Skeleton className="h-3 w-24" />
        </div>
        <Skeleton className="h-4 w-full mb-3" />
        <Skeleton className="h-16 w-full" />
      </div>
    </div>
  );
}

function AddActivityDialog({ 
  open, 
  onOpenChange,
  onSubmit 
}: { 
  open: boolean; 
  onOpenChange: (open: boolean) => void;
  onSubmit: (activity: Partial<Activity>) => void;
}) {
  const [activity, setActivity] = useState<Partial<Activity>>({
    type: 'email',
    title: '',
    description: '',
    content: ''
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(activity);
    onOpenChange(false);
    setActivity({
      type: 'email',
      title: '',
      description: '',
      content: ''
    });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle>Add New Activity</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Type</label>
            <Select
              value={activity.type}
              onValueChange={(value) => setActivity({ ...activity, type: value as Activity['type'] })}
            >
              <SelectTrigger>
                <SelectValue placeholder="Select activity type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="email">Email</SelectItem>
                <SelectItem value="meeting">Meeting</SelectItem>
                <SelectItem value="note">Note</SelectItem>
              </SelectContent>
            </Select>
          </div>
          
          <div className="space-y-2">
            <label className="text-sm font-medium">Title</label>
            <Input
              value={activity.title}
              onChange={(e) => setActivity({ ...activity, title: e.target.value })}
              placeholder={activity.type === 'note' ? 'Note title' : 'Activity title'}
              required
            />
          </div>
          
          <div className="space-y-2">
            <label className="text-sm font-medium">{activity.type === 'note' ? 'Note Content' : 'Description'}</label>
            {activity.type === 'note' ? (
              <Textarea
                value={activity.content}
                onChange={(e) => setActivity({ ...activity, content: e.target.value, description: e.target.value.slice(0, 100) + '...' })}
                placeholder="Write your note here"
                className="min-h-[200px]"
                required
              />
            ) : (
              <Input
                value={activity.description}
                onChange={(e) => setActivity({ ...activity, description: e.target.value })}
                placeholder="Brief description"
                required
              />
            )}
          </div>
          
          {activity.type !== 'note' && (
            <div className="space-y-2">
              <label className="text-sm font-medium">Content (optional)</label>
              <Textarea
                value={activity.content}
                onChange={(e) => setActivity({ ...activity, content: e.target.value })}
                placeholder="Additional details or content"
              />
            </div>
          )}

          <div className="flex justify-end gap-2">
            <Button variant="outline" type="button" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit">Add {activity.type === 'note' ? 'Note' : 'Activity'}</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}

function NoteDialog({ 
  open, 
  onOpenChange,
  onSubmit,
  initialNote
}: { 
  open: boolean; 
  onOpenChange: (open: boolean) => void;
  onSubmit: (note: Partial<Note>) => void;
  initialNote?: Note;
}) {
  const [note, setNote] = useState<Partial<Note>>(() => initialNote || {
    title: '',
    content: ''
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(note);
    onOpenChange(false);
    if (!initialNote) {
      setNote({
        title: '',
        content: ''
      });
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{initialNote ? 'Edit Note' : 'Add New Note'}</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Title</label>
            <Input
              value={note.title}
              onChange={(e) => setNote({ ...note, title: e.target.value })}
              placeholder="Note title"
              required
            />
          </div>
          
          <div className="space-y-2">
            <label className="text-sm font-medium">Content</label>
            <Textarea
              value={note.content}
              onChange={(e) => setNote({ ...note, content: e.target.value })}
              placeholder="Note content"
              required
              className="min-h-[200px]"
            />
          </div>

          <div className="flex justify-end gap-2">
            <Button variant="outline" type="button" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit">{initialNote ? 'Save Changes' : 'Add Note'}</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}

function DocumentIcon({ type }: { type: string }) {
  switch (type) {
    case 'pdf':
      return <RiFileTextLine className="w-5 h-5 text-red-500" />;
    case 'excel':
      return <RiFileExcelLine className="w-5 h-5 text-green-500" />;
    case 'image':
      return <RiImageLine className="w-5 h-5 text-blue-500" />;
    default:
      return <RiFile3Line className="w-5 h-5 text-gray-500" />;
  }
}

function formatFileSize(bytes: number): string {
  const units = ['B', 'KB', 'MB', 'GB'];
  let size = bytes;
  let unitIndex = 0;
  
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024;
    unitIndex++;
  }
  
  return `${size.toFixed(1)} ${units[unitIndex]}`;
}

function truncateText(text: string, maxLength: number) {
  if (text.length <= maxLength) return text;
  return text.slice(0, maxLength - 3) + '...';
}

function UploadDocumentModal({ 
  onUpload 
}: { 
  onUpload: (document: Document) => void;
}) {
  const [open, setOpen] = useState(false);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [title, setTitle] = useState('');
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setSelectedFile(file);
      // Set the title to the file name without the extension
      const fileName = file.name.replace(/\.[^/.]+$/, "");
      setTitle(fileName);
    }
  };

  const handleDrop = (e: React.DragEvent<HTMLLabelElement>) => {
    e.preventDefault();
    const file = e.dataTransfer.files?.[0];
    if (file) {
      setSelectedFile(file);
      const fileName = file.name.replace(/\.[^/.]+$/, "");
      setTitle(fileName);
    }
  };

  const handleDragOver = (e: React.DragEvent<HTMLLabelElement>) => {
    e.preventDefault();
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedFile) return;

    const newDocument: Document = {
      id: Math.random().toString(36).substring(7),
      name: title || selectedFile.name,
      type: selectedFile.type.includes('pdf') ? 'pdf' 
        : selectedFile.type.includes('excel') || selectedFile.type.includes('spreadsheet') ? 'excel'
        : selectedFile.type.includes('image') ? 'image'
        : 'other',
      size: selectedFile.size,
      uploadedAt: new Date(),
      uploadedBy: 'Current User'
    };

    onUpload(newDocument);
    setOpen(false);
    setSelectedFile(null);
    setTitle('');
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button
          variant="outline"
          size="sm"
          className="flex items-center gap-2"
        >
          <RiUploadLine className="w-4 h-4" />
          Upload Document
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Upload Document</DialogTitle>
          <DialogDescription>
            Upload a document to attach to this investor.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="mt-6 space-y-6">
          <div className="space-y-2">
            <Label htmlFor="title">Document Title</Label>
            <Input
              id="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Enter document title"
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="document">Document File</Label>
            <div className="flex items-center justify-center w-full">
              <label
                htmlFor="document"
                className="flex flex-col items-center justify-center w-full h-32 border-2 border-dashed rounded-lg cursor-pointer border-input bg-background/50 hover:bg-hover"
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
                      <p className="text-xs text-muted-foreground">PDF, Excel, or image files</p>
                    </>
                  )}
                </div>
                <Input
                  ref={fileInputRef}
                  id="document"
                  type="file"
                  className="hidden"
                  accept=".pdf,.xlsx,.xls,.doc,.docx,.png,.jpg,.jpeg"
                  onChange={handleFileChange}
                />
              </label>
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="ghost"
              onClick={() => {
                setOpen(false);
                setSelectedFile(null);
                setTitle('');
              }}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={!selectedFile || !title}
            >
              Upload
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

function DocumentsTab({ isReadOnly }: { isReadOnly: boolean }) {
  const [documents, setDocuments] = useState<Document[]>(mockDocuments);
  const [documentToDelete, setDocumentToDelete] = useState<Document | null>(null);

  const handleUpload = (document: Document) => {
    setDocuments(prev => [document, ...prev]);
    toast.success('Document uploaded successfully');
  };

  const handleDownload = (document: Document) => {
    // In a real implementation, you would download the file from your backend here
    toast.success(`Downloading ${document.name}`);
  };

  const handleDelete = (document: Document) => {
    setDocumentToDelete(document);
  };

  const confirmDelete = () => {
    if (!documentToDelete) return;
    
    setDocuments(prev => prev.filter(doc => doc.id !== documentToDelete.id));
    toast.success('Document deleted successfully');
    setDocumentToDelete(null);
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold">Documents</h3>
        {!isReadOnly && (
          <UploadDocumentModal onUpload={handleUpload} />
        )}
      </div>

      {documents.length === 0 ? (
        <div className="text-center py-8 text-muted-foreground">
          No documents uploaded yet
        </div>
      ) : (
        <div className="space-y-2">
          {documents.map(document => (
            <div
              key={document.id}
              className="flex items-center justify-between p-3 rounded-lg border bg-card"
            >
              <div className="flex items-center gap-3">
                <DocumentIcon type={document.type} />
                <div>
                  <p className="font-medium">{document.name}</p>
                  <p className="text-sm text-muted-foreground">
                    {formatFileSize(document.size)} • Uploaded by {document.uploadedBy} on{' '}
                    {new Date(document.uploadedAt).toLocaleDateString()}
                  </p>
                </div>
              </div>
              <div className="flex items-center gap-2">
                {!isReadOnly && (
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleDownload(document)}
                  >
                    <RiDownloadLine className="w-4 h-4" />
                  </Button>
                )}
                {!isReadOnly && (
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleDelete(document)}
                  >
                    <RiDeleteBinLine className="w-4 h-4" />
                  </Button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}

      <AlertDialog open={!!documentToDelete} onOpenChange={() => setDocumentToDelete(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Document</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete "{documentToDelete?.name}"? This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            {!isReadOnly && (
              <AlertDialogAction onClick={confirmDelete}>Delete</AlertDialogAction>
            )}
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}

function EditInvestorDialog({ 
  open, 
  onOpenChange,
  investor,
  onSave
}: { 
  open: boolean; 
  onOpenChange: (open: boolean) => void;
  investor: Card | null;
  onSave: (updatedInvestor: Card) => void;
}) {
  const [editedInvestor, setEditedInvestor] = useState<Card | null>(null);
  const [hoveredRating, setHoveredRating] = useState<number | null>(null);

  useEffect(() => {
    if (investor) {
      setEditedInvestor({ ...investor });
    }
  }, [investor]);

  if (!editedInvestor) return null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (editedInvestor) {
      onSave(editedInvestor);
      onOpenChange(false);
    }
  };

  const displayRating = hoveredRating !== null ? hoveredRating : editedInvestor.rating;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Investor Details</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label>Amount</Label>
            <NumericFormat
              value={editedInvestor.amount.replace(/[^0-9.]/g, '')}
              onValueChange={(values) => {
                const { value } = values;
                // Format the value with 'M' suffix
                const formattedValue = `${value}M`;
                setEditedInvestor({ ...editedInvestor, amount: formattedValue });
              }}
              thousandSeparator
              prefix="$"
              customInput={Input}
              placeholder="Enter amount (e.g. 5M)"
              className="[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
            />
          </div>

          <div className="space-y-2">
            <Label>Check Size</Label>
            <NumericFormat
              value={editedInvestor.checkSize.replace(/[^0-9.]/g, '')}
              onValueChange={(values) => {
                const { value } = values;
                // Format the value with 'M' suffix
                const formattedValue = `${value}M`;
                setEditedInvestor({ ...editedInvestor, checkSize: formattedValue });
              }}
              thousandSeparator
              prefix="$"
              customInput={Input}
              placeholder="Enter check size"
              className="[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
            />
          </div>
          
          <div className="space-y-2">
            <Label>Initial Contact Date</Label>
            <Input
              type="date"
              value={editedInvestor.initialContactDate}
              onChange={(e) => setEditedInvestor({ ...editedInvestor, initialContactDate: e.target.value })}
            />
          </div>
          
          <div className="space-y-2">
            <Label>Lead Investor</Label>
            <div className="flex items-center space-x-2">
              <Switch
                checked={editedInvestor.isLeadInvestor}
                onCheckedChange={(checked) => setEditedInvestor({ ...editedInvestor, isLeadInvestor: checked })}
              />
              <span className="text-sm text-muted-foreground">
                {editedInvestor.isLeadInvestor ? 'Yes' : 'No'}
              </span>
            </div>
          </div>
          
          <div className="space-y-2">
            <Label>Rating</Label>
            <div className="flex items-center gap-1">
              {[1, 2, 3, 4, 5].map((star) => (
                <Button
                  key={star}
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="hover:bg-transparent"
                  onClick={() => setEditedInvestor({ ...editedInvestor, rating: star })}
                  onMouseEnter={() => setHoveredRating(star)}
                  onMouseLeave={() => setHoveredRating(null)}
                >
                  {star <= displayRating ? (
                    <RiStarFill className="w-6 h-6 text-yellow-400" />
                  ) : (
                    <RiStarLine className="w-6 h-6 text-muted-foreground" />
                  )}
                </Button>
              ))}
              <span className="ml-2 text-sm text-muted-foreground">
                {displayRating} of 5
              </span>
            </div>
          </div>

          <div className="flex justify-end gap-2">
            <Button variant="outline" type="button" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit">Save Changes</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}

export function InvestorDetailsDrawer({
  open,
  onOpenChange,
  investor,
  isArchived = false,
}: InvestorDetailsDrawerProps) {
  const [activeTab, setActiveTab] = useState("overview");
  const [activities, setActivities] = useState<Activity[]>([]);
  const [isAddingActivity, setIsAddingActivity] = useState(false);
  const [isEditingInvestor, setIsEditingInvestor] = useState(false);
  
  // Infinite scroll states
  const [page, setPage] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const observerTarget = useRef<HTMLDivElement>(null);

  const loadMoreActivities = async () => {
    if (!investor || isLoading || !hasMore) return;
    
    setIsLoading(true);
    try {
      const newActivities = await fetchActivities(page, investor.id);
      if (newActivities.length === 0 || activities.length + newActivities.length >= TOTAL_ACTIVITIES_LIMIT) {
        setHasMore(false);
      }
      setActivities(prev => [...prev, ...newActivities]);
      setPage(prev => prev + 1);
    } catch (error) {
      console.error('Error loading activities:', error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    const observer = new IntersectionObserver(
      entries => {
        if (entries[0].isIntersecting && hasMore && !isLoading) {
          loadMoreActivities();
        }
      },
      { threshold: 0.1 }
    );

    if (observerTarget.current) {
      observer.observe(observerTarget.current);
    }

    return () => observer.disconnect();
  }, [hasMore, isLoading, page]);

  // Reset states when investor changes
  useEffect(() => {
    if (investor) {
      setActivities([]);
      setPage(0);
      setHasMore(true);
      loadMoreActivities();
    }
  }, [investor?.id]);

  const handleAddActivity = (newActivity: Partial<Activity>) => {
    const activity: Activity = {
      ...newActivity,
      id: Math.random().toString(36).substr(2, 9),
      timestamp: new Date().toISOString(),
      type: newActivity.type as Activity['type'],
      title: newActivity.title || '',
      description: newActivity.description || ''
    };
    
    setActivities(prev => [activity, ...prev]);
  };

  const getActivityIcon = (type: Activity['type']) => {
    switch (type) {
      case 'email':
        return <RiMailLine className="w-4 h-4 text-primary" />;
      case 'meeting':
        return <RiCalendarLine className="w-4 h-4 text-primary" />;
      case 'note':
        return <RiFileTextLine className="w-4 h-4 text-primary" />;
      default:
        return <RiMailLine className="w-4 h-4 text-primary" />;
    }
  };

  const handleSaveInvestor = (updatedInvestor: Card) => {
    // In a real implementation, you would update the investor in your backend here
    onOpenChange(false);
  };

  if (!investor) return null;

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent 
        className="w-full max-w-[90%] lg:max-w-[75%] 2xl:max-w-[1400px] overflow-y-auto" 
      >
        <div className="h-full flex flex-col">
          <SheetHeader className="flex-none">
            <div className="flex items-center justify-between">
              <SheetTitle className="text-xl">{investor.title}</SheetTitle>
              <div className="flex items-center gap-2">
                {!isArchived && (
                  <Button variant="ghost" size="icon" onClick={() => setIsAddingActivity(true)}>
                    <RiAddLine className="w-4 h-4" />
                  </Button>
                )}
                <Button variant="ghost" size="icon" onClick={() => onOpenChange(false)}>
                  <RiCloseLine className="w-4 h-4" />
                </Button>
              </div>
            </div>
            <div className="flex items-center justify-between">
              <div className="space-x-2">
                <Badge variant="outline">{investor.stage}</Badge>
                <Badge variant="secondary" className="bg-primary/10 text-primary">
                  ${investor.amount}
                </Badge>
                {isArchived && (
                  <Badge variant="outline" className="text-muted-foreground">
                    <RiArchiveFill className="w-3 h-3 mr-1 inline" />
                    Archived
                  </Badge>
                )}
              </div>
            </div>
          </SheetHeader>

          <div className="mt-6">
            <Tabs value={activeTab} onValueChange={setActiveTab}>
              <TabsList className="w-full">
                <TabsTrigger value="overview" className="flex-1">Overview</TabsTrigger>
                <TabsTrigger value="activity" className="flex-1">Activity</TabsTrigger>
                <TabsTrigger value="documents" className="flex-1">Documents</TabsTrigger>
              </TabsList>

              <TabsContent value="overview" className="mt-6">
                <div className="space-y-6">
                  <div className="bg-card rounded-lg p-4 border">
                    <h3 className="font-medium mb-3">Contact Information</h3>
                    <div className="space-y-3">
                      <div className="flex items-center">
                        <Avatar className="w-10 h-10">
                          <AvatarImage
                            src={investor.contact.image}
                            alt={investor.contact.name}
                          />
                          <AvatarFallback>
                            {investor.contact.name
                              .split(" ")
                              .map((n) => n[0])
                              .join("")}
                          </AvatarFallback>
                        </Avatar>
                        <div className="ml-3">
                          <p className="font-medium">{investor.contact.name}</p>
                          <p className="text-sm text-muted-foreground">
                            Primary Contact
                          </p>
                        </div>
                      </div>
                      <div className="flex items-center text-sm">
                        <RiMailLine className="w-4 h-4 mr-2" />
                        <span>contact@example.com</span>
                      </div>
                      <div className="flex items-center text-sm">
                        <RiPhoneLine className="w-4 h-4 mr-2" />
                        <span>+1 (555) 123-4567</span>
                      </div>
                    </div>
                  </div>

                  <div className="bg-card rounded-lg p-4 border">
                    <h3 className="font-medium mb-3">Deal Information</h3>
                    <div className="space-y-3 text-sm">
                      <div className="flex items-center justify-between">
                        <span className="text-muted-foreground">Deal Size</span>
                        <span className="font-medium">${investor.amount}</span>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-muted-foreground">Check Size</span>
                        <span className="font-medium">{investor.checkSize}</span>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-muted-foreground">Initial Contact</span>
                        <span className="font-medium">{investor.initialContactDate}</span>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-muted-foreground">Lead Investor</span>
                        <Badge variant={investor.isLeadInvestor ? "default" : "outline"}>
                          {investor.isLeadInvestor ? "Yes" : "No"}
                        </Badge>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-muted-foreground">Rating</span>
                        <div className="flex items-center">
                          {[1, 2, 3, 4, 5].map((star) => (
                            <RiStarLine
                              key={star}
                              className={`w-4 h-4 ${star <= (investor.rating || 0) ? 'text-yellow-400' : 'text-muted-foreground'}`}
                            />
                          ))}
                        </div>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-muted-foreground">Stage</span>
                        <Badge variant="outline">{investor.stage}</Badge>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-muted-foreground">Due Date</span>
                        <span className="font-medium">{investor.dueDate}</span>
                      </div>
                    </div>
                  </div>

                  <div className="bg-card rounded-lg p-4 border">
                    <h3 className="font-medium mb-3">Suggested Actions</h3>
                    <div className="space-y-2">
                      <Button 
                        variant="outline" 
                        className="w-full justify-start"
                        onClick={() => setActiveTab("activity")}
                      >
                        <RiCalendarLine className="w-4 h-4 mr-2" />
                        Add Activity or Note
                      </Button>
                      <Button 
                        variant="outline" 
                        className="w-full justify-start"
                        onClick={() => setActiveTab("documents")}
                      >
                        <RiFileTextLine className="w-4 h-4 mr-2" />
                        Upload Documents
                      </Button>
                    </div>
                  </div>
                </div>
              </TabsContent>

              <TabsContent value="activity">
                <div className="space-y-6">
                  <div className="flex items-center justify-between">
                    <div className="text-sm text-muted-foreground">
                      Showing {activities.length} of {activities.length >= TOTAL_ACTIVITIES_LIMIT ? 'maximum ' : ''}{TOTAL_ACTIVITIES_LIMIT} activities
                    </div>
                    {!isArchived && (
                      <Button 
                        onClick={() => setIsAddingActivity(true)}
                        size="sm"
                        className="gap-2"
                      >
                        <RiAddLine className="w-4 h-4" />
                        Add Activity or Note
                      </Button>
                    )}
                  </div>
                  
                  <div className="relative space-y-6 pl-8 before:absolute before:left-3 before:top-2 before:bottom-0 before:w-[2px] before:bg-border">
                    {activities.map((activity) => (
                      <div key={activity.id} className="relative">
                        <div className="absolute -left-[27px] bg-background p-1 border rounded-full">
                          {getActivityIcon(activity.type)}
                        </div>
                        <div className="bg-card rounded-lg p-4 border">
                          <div className="flex items-center justify-between mb-2">
                            <h4 className="font-medium">{activity.title}</h4>
                            <span className="text-xs text-muted-foreground">
                              {new Date(activity.timestamp).toLocaleString()}
                            </span>
                          </div>
                          <p className="text-sm text-muted-foreground mb-3">
                            {activity.description}
                          </p>
                          {activity.content && (
                            <div className="bg-muted/50 rounded-md p-3 text-sm">
                              <p className="whitespace-pre-wrap">{activity.content}</p>
                            </div>
                          )}
                        </div>
                      </div>
                    ))}

                    {/* Loading state */}
                    {isLoading && (
                      <>
                        <ActivitySkeleton />
                        <ActivitySkeleton />
                        <ActivitySkeleton />
                      </>
                    )}

                    {/* Intersection observer target */}
                    <div ref={observerTarget} className="h-4" />

                    {/* End of list message */}
                    {!hasMore && activities.length > 0 && (
                      <div className="text-center text-sm text-muted-foreground py-4">
                        {activities.length >= TOTAL_ACTIVITIES_LIMIT 
                          ? `Maximum limit of ${TOTAL_ACTIVITIES_LIMIT} activities reached`
                          : "No more activities to load"}
                      </div>
                    )}
                  </div>
                </div>
              </TabsContent>

              <TabsContent value="documents">
                <DocumentsTab isReadOnly={isArchived} />
              </TabsContent>
            </Tabs>
          </div>

          {!isArchived && (
            <>
              <EditInvestorDialog
                open={isEditingInvestor}
                onOpenChange={setIsEditingInvestor}
                investor={investor}
                onSave={handleSaveInvestor}
              />

              <AddActivityDialog
                open={isAddingActivity}
                onOpenChange={setIsAddingActivity}
                onSubmit={handleAddActivity}
              />
            </>
          )}
        </div>
      </SheetContent>
    </Sheet>
  );
} 