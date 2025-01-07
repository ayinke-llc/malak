import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Mail,
  Phone,
  MapPin,
  Building2,
  Calendar,
  Pencil,
  Trash2,
  Users,
  BarChart3,
  Clock,
  FileText,
  LayoutDashboard,
  FolderOpen,
  ChevronRight
} from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Bar, BarChart, ResponsiveContainer, XAxis, YAxis } from "recharts";
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
import { toast } from "sonner";
import CreateContactModal from "@/components/ui/contacts/modal";

type TimePeriod = 'days' | 'weeks' | 'months';

interface ContactDetailsProps {
  reference: string;
}

const ContactDetails = ({ reference }: ContactDetailsProps) => {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [timePeriod, setTimePeriod] = useState<TimePeriod>('months');
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);

  // TODO: Fetch contact details using the reference

  const getChartData = (period: TimePeriod) => {
    switch (period) {
      case 'days':
        return [
          { name: "Mon", total: 2 },
          { name: "Tue", total: 5 },
          { name: "Wed", total: 7 },
          { name: "Thu", total: 3 },
          { name: "Fri", total: 8 },
          { name: "Sat", total: 4 },
          { name: "Sun", total: 6 },
        ];
      case 'weeks':
        return [
          { name: "Week 1", total: 12 },
          { name: "Week 2", total: 8 },
          { name: "Week 3", total: 15 },
          { name: "Week 4", total: 10 },
        ];
      case 'months':
      default:
        return [
          { name: "May", total: 0 },
          { name: "Jul", total: 0 },
          { name: "Sep", total: 5 },
          { name: "Nov", total: 15 },
          { name: "Jan 2024", total: 10 },
          { name: "Mar", total: 10 },
        ];
    }
  };

  const chartdata = getChartData(timePeriod);

  const handleDelete = async () => {
    try {
      setIsLoading(true);
      // TODO: API call to delete contact
      await fetch(`/api/contacts/${reference}`, {
        method: 'DELETE',
      });
      toast.success("Contact deleted successfully");
      router.push('/contacts');
    } catch (error) {
      toast.error("Failed to delete contact");
    } finally {
      setIsLoading(false);
      setShowDeleteDialog(false);
    }
  };

  return (
    <div className="mt-6 space-y-6">
      <Card className="shadow-sm">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-7">
          <div className="space-y-1">
            <CardTitle className="text-2xl font-bold tracking-tight text-foreground">John Doe</CardTitle>
            <CardDescription className="text-base text-muted-foreground">Contact Information</CardDescription>
          </div>
          <div className="flex gap-3">
            <Button
              variant="outline"
              size="icon"
              onClick={() => setShowEditModal(true)}
              className="h-9 w-9"
            >
              <Pencil className="h-4 w-4" />
            </Button>
            <Button
              variant="destructive"
              size="icon"
              onClick={() => setShowDeleteDialog(true)}
              disabled={isLoading}
              className="h-9 w-9"
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <Tabs defaultValue="details" className="w-full">
            <TabsList className="mb-4">
              <TabsTrigger value="details" className="text-sm">Details</TabsTrigger>
              <TabsTrigger value="activity" className="text-sm">Activity</TabsTrigger>
              <TabsTrigger value="notes" className="text-sm">Notes</TabsTrigger>
            </TabsList>

            <TabsContent value="details">
              <div className="grid gap-8 mt-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                  <div className="space-y-6">
                    <div className="flex items-start gap-3 p-4 rounded-lg bg-muted/50">
                      <Mail className="h-5 w-5 text-muted-foreground mt-0.5" />
                      <div>
                        <p className="text-sm font-semibold text-foreground mb-1">Email</p>
                        <p className="text-sm text-muted-foreground">john.doe@example.com</p>
                      </div>
                    </div>
                    <div className="flex items-start gap-3 p-4 rounded-lg bg-muted/50">
                      <Phone className="h-5 w-5 text-muted-foreground mt-0.5" />
                      <div>
                        <p className="text-sm font-semibold text-foreground mb-1">Phone</p>
                        <p className="text-sm text-muted-foreground">+1 (555) 123-4567</p>
                      </div>
                    </div>
                  </div>
                  <div className="space-y-6">
                    <div className="flex items-start gap-3 p-4 rounded-lg bg-muted/50">
                      <Building2 className="h-5 w-5 text-muted-foreground mt-0.5" />
                      <div>
                        <p className="text-sm font-semibold text-foreground mb-1">Company</p>
                        <p className="text-sm text-muted-foreground">Acme Inc.</p>
                      </div>
                    </div>
                    <div className="flex items-start gap-3 p-4 rounded-lg bg-muted/50">
                      <MapPin className="h-5 w-5 text-muted-foreground mt-0.5" />
                      <div>
                        <p className="text-sm font-semibold text-foreground mb-1">Address</p>
                        <p className="text-sm text-muted-foreground">
                          123 Business Street, Suite 100<br />
                          San Francisco, CA 94105
                        </p>
                      </div>
                    </div>
                  </div>
                </div>

                <div className="border-t pt-6">
                  <h4 className="text-sm font-semibold text-foreground mb-4">Contact Lists</h4>
                  <div className="flex items-start gap-3">
                    <Users className="h-5 w-5 text-muted-foreground mt-0.5" />
                    <div className="flex-1">
                      <div className="flex flex-wrap gap-2">
                        <Badge variant="secondary" className="px-3 py-1">Customers</Badge>
                        <Badge variant="secondary" className="px-3 py-1">VIP</Badge>
                        <Badge variant="secondary" className="px-3 py-1">Newsletter</Badge>
                      </div>
                    </div>
                  </div>
                </div>

                <div className="border-t pt-6">
                  <h4 className="text-sm font-semibold text-foreground mb-4">Additional Information</h4>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className="flex items-start gap-3 p-4 rounded-lg bg-muted/50">
                      <Calendar className="h-5 w-5 text-muted-foreground mt-0.5" />
                      <div>
                        <p className="text-sm font-semibold text-foreground mb-1">Created</p>
                        <p className="text-sm text-muted-foreground">Jan 1, 2024</p>
                      </div>
                    </div>
                    <div className="flex items-start gap-3 p-4 rounded-lg bg-muted/50">
                      <Calendar className="h-5 w-5 text-muted-foreground mt-0.5" />
                      <div>
                        <p className="text-sm font-semibold text-foreground mb-1">Last Updated</p>
                        <p className="text-sm text-muted-foreground">Jan 15, 2024</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </TabsContent>

            <TabsContent value="activity">
              <div className="space-y-10 mt-6">
                {/* Engagement Trend Section */}
                <div>
                  <div className="flex items-center justify-between mb-6">
                    <h3 className="text-lg font-semibold text-foreground flex items-center gap-2">
                      <BarChart3 className="h-5 w-5 text-muted-foreground" />
                      Engagement trend
                    </h3>
                    <div className="inline-flex items-center rounded-md bg-muted p-1 text-muted-foreground">
                      <Button
                        variant={timePeriod === 'days' ? 'secondary' : 'ghost'}
                        size="sm"
                        className="text-xs px-3 py-1.5"
                        onClick={() => setTimePeriod('days')}
                      >
                        Days
                      </Button>
                      <Button
                        variant={timePeriod === 'weeks' ? 'secondary' : 'ghost'}
                        size="sm"
                        className="text-xs px-3 py-1.5"
                        onClick={() => setTimePeriod('weeks')}
                      >
                        Weeks
                      </Button>
                      <Button
                        variant={timePeriod === 'months' ? 'secondary' : 'ghost'}
                        size="sm"
                        className="text-xs px-3 py-1.5"
                        onClick={() => setTimePeriod('months')}
                      >
                        Months
                      </Button>
                    </div>
                  </div>
                  <Card className="shadow-sm">
                    <CardContent className="pl-2 pt-6">
                      <div className="h-[250px] w-full">
                        <ResponsiveContainer width="100%" height="100%">
                          <BarChart data={chartdata}>
                            <XAxis
                              dataKey="name"
                              stroke="currentColor"
                              fontSize={12}
                              tickLine={false}
                              axisLine={false}
                              className="text-muted-foreground"
                            />
                            <YAxis
                              stroke="currentColor"
                              fontSize={12}
                              tickLine={false}
                              axisLine={false}
                              tickFormatter={(value) => `${value}`}
                              className="text-muted-foreground"
                            />
                            <Bar
                              dataKey="total"
                              fill="currentColor"
                              radius={[4, 4, 0, 0]}
                              className="fill-primary"
                            />
                          </BarChart>
                        </ResponsiveContainer>
                      </div>
                    </CardContent>
                  </Card>
                </div>

                {/* Recent Activity Section */}
                <div>
                  <h3 className="text-lg font-semibold text-foreground mb-6 flex items-center gap-2">
                    <Clock className="h-5 w-5 text-muted-foreground" />
                    Recent activity
                  </h3>
                  <div className="space-y-4">
                    <div className="flex items-start gap-4 p-5 rounded-lg border bg-card">
                      <div className="h-10 w-10 rounded-full bg-muted flex items-center justify-center">
                        <Users className="h-5 w-5 text-muted-foreground" />
                      </div>
                      <div className="flex-1">
                        <div className="flex items-center justify-between mb-1">
                          <p className="text-sm font-semibold text-foreground">Added to VIP list</p>
                          <span className="text-sm text-muted-foreground">2 days ago</span>
                        </div>
                        <p className="text-sm text-muted-foreground">Contact was added to the VIP contact list</p>
                      </div>
                    </div>

                    <div className="flex items-start gap-4 p-5 rounded-lg border bg-card">
                      <div className="h-10 w-10 rounded-full bg-muted flex items-center justify-center">
                        <Mail className="h-5 w-5 text-muted-foreground" />
                      </div>
                      <div className="flex-1">
                        <div className="flex items-center justify-between mb-1">
                          <p className="text-sm font-semibold text-foreground">Email sent</p>
                          <span className="text-sm text-muted-foreground">5 days ago</span>
                        </div>
                        <p className="text-sm text-muted-foreground">Monthly newsletter was sent to contact</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </TabsContent>

            <TabsContent value="notes">
              <div className="mt-6">
                <p className="text-sm text-muted-foreground">No notes available</p>
              </div>
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>

      {/* Latest Shared Items Section */}
      <div className="space-y-8 mt-10">
        <h2 className="text-2xl font-semibold tracking-tight text-foreground">Latest shared items</h2>

        {/* Updates Section */}
        <div>
          <div className="flex items-center mb-4">
            <div className="flex items-center gap-2">
              <FileText className="h-5 w-5 text-muted-foreground" />
              <h3 className="text-lg font-semibold text-foreground">Updates</h3>
            </div>
          </div>
          <div className="space-y-2">
            <div className="flex items-center justify-between p-3 -mx-2 rounded-md transition-colors hover:bg-accent hover:text-accent-foreground cursor-pointer">
              <span className="text-sm font-medium text-foreground">December investor update</span>
              <span className="text-sm text-muted-foreground">Sent 2 months ago</span>
            </div>
            <div className="flex items-center justify-between p-3 -mx-2 rounded-md transition-colors hover:bg-accent hover:text-accent-foreground cursor-pointer">
              <span className="text-sm font-medium text-foreground">Duplicate of Test update</span>
              <span className="text-sm text-muted-foreground">Sent 3 months ago</span>
            </div>
          </div>
        </div>

        {/* Dashboards Section */}
        <div>
          <div className="flex items-center mb-4">
            <div className="flex items-center gap-2">
              <LayoutDashboard className="h-5 w-5 text-muted-foreground" />
              <h3 className="text-lg font-semibold text-foreground">Dashboards</h3>
            </div>
          </div>
          <div className="space-y-2">
            <div className="flex items-center justify-between p-3 -mx-2 rounded-md transition-colors hover:bg-accent hover:text-accent-foreground cursor-pointer">
              <span className="text-sm font-medium text-foreground">Financial metrics</span>
              <span className="text-sm text-muted-foreground">Shared 17 hours ago</span>
            </div>
          </div>
        </div>

        {/* Data Rooms Section */}
        <div>
          <div className="flex items-center mb-4">
            <div className="flex items-center gap-2">
              <FolderOpen className="h-5 w-5 text-muted-foreground" />
              <h3 className="text-lg font-semibold text-foreground">Data rooms</h3>
            </div>
          </div>
          <div className="space-y-2">
            <div className="flex items-center justify-between p-3 -mx-2 rounded-md transition-colors hover:bg-accent hover:text-accent-foreground cursor-pointer">
              <span className="text-sm font-medium text-foreground">Data room example</span>
              <span className="text-sm text-muted-foreground">Shared 4 months ago</span>
            </div>
          </div>
        </div>
      </div>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. This will permanently delete the contact
              and remove their data from our servers.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={isLoading}>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDelete}
              disabled={isLoading}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {isLoading ? "Deleting..." : "Delete"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
};

export default ContactDetails;
