import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Textarea } from "@/components/ui/textarea"
import { CheckCircle2, Clock, Plus, XCircle } from "lucide-react"
import Link from "next/link"

type FundingStage = 
  | "pre-seed"
  | "seed"
  | "series-a"
  | "series-b"
  | "series-c"
  | "series-d"
  | "growth"
  | "bridge"
  | "ipo"

const FUNDING_STAGES: { value: FundingStage; label: string; description: string }[] = [
  {
    value: "pre-seed",
    label: "Pre-Seed",
    description: "Very early stage funding to develop initial product"
  },
  {
    value: "seed",
    label: "Seed",
    description: "Early stage funding to validate product market fit"
  },
  {
    value: "series-a",
    label: "Series A",
    description: "First significant round of venture capital financing"
  },
  {
    value: "series-b",
    label: "Series B",
    description: "Funding for business development and market growth"
  },
  {
    value: "series-c",
    label: "Series C",
    description: "Scaling, expansion, and possible acquisitions"
  },
  {
    value: "series-d",
    label: "Series D",
    description: "Additional scaling or preparation for exit"
  },
  {
    value: "growth",
    label: "Growth",
    description: "Late-stage funding for rapid expansion"
  },
  {
    value: "bridge",
    label: "Bridge",
    description: "Short-term funding between major rounds"
  },
  {
    value: "ipo",
    label: "IPO",
    description: "Preparation for public offering"
  }
]

// Example data - replace with actual data from your API
const exampleRounds = [
  {
    id: 1,
    name: "2024 Growth Round",
    stage: "series-a",
    target: 5000000,
    raised: 2500000,
    status: "active",
    description: "Growth and expansion funding round",
    deadline: "2024-12-31",
    closedReason: null
  },
  {
    id: 2,
    name: "Initial Funding",
    stage: "seed",
    target: 1000000,
    raised: 1000000,
    status: "successful",
    description: "Initial seed funding for product development",
    deadline: "2024-06-30",
    closedReason: "Fully funded"
  },
  {
    id: 3,
    name: "Q3 Bridge",
    stage: "bridge",
    target: 2000000,
    raised: 100000,
    status: "cancelled",
    description: "Bridge funding for market expansion",
    deadline: "2024-09-30",
    closedReason: "Market conditions unfavorable"
  },
  {
    id: 4,
    name: "Scale-up 2025",
    stage: "series-b",
    target: 10000000,
    raised: 0,
    status: "upcoming",
    description: "Scale operations and enter new markets",
    deadline: "2025-03-31",
    closedReason: null
  },
  {
    id: 5,
    name: "Strategic Round",
    stage: "growth",
    target: 3000000,
    raised: 3000000,
    status: "successful",
    description: "Partnership and technology development",
    deadline: "2024-11-30",
    closedReason: "Target reached"
  },
  {
    id: 6,
    name: "Early Stage",
    stage: "pre-seed",
    target: 500000,
    raised: 200000,
    status: "cancelled",
    description: "Initial angel investment round",
    deadline: "2024-08-15",
    closedReason: "Pivoting business model"
  }
]

type RoundStatus = "active" | "successful" | "cancelled" | "upcoming"

function formatCurrency(amount: number) {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    notation: 'compact',
    maximumFractionDigits: 1,
  }).format(amount)
}

function getStatusConfig(status: RoundStatus) {
  switch (status) {
    case 'successful':
      return {
        color: 'bg-green-500',
        textColor: 'text-green-700',
        bgColor: 'bg-green-50',
        icon: CheckCircle2,
        label: 'Successful'
      }
    case 'active':
      return {
        color: 'bg-blue-500',
        textColor: 'text-blue-700',
        bgColor: 'bg-blue-50',
        icon: Clock,
        label: 'Active'
      }
    case 'upcoming':
      return {
        color: 'bg-yellow-500',
        textColor: 'text-yellow-700',
        bgColor: 'bg-yellow-50',
        icon: Clock,
        label: 'Upcoming'
      }
    case 'cancelled':
      return {
        color: 'bg-gray-500',
        textColor: 'text-gray-700',
        bgColor: 'bg-gray-50',
        icon: XCircle,
        label: 'Cancelled'
      }
    default:
      return {
        color: 'bg-gray-500',
        textColor: 'text-gray-700',
        bgColor: 'bg-gray-50',
        icon: Clock,
        label: status
      }
  }
}

export default function FundraisingBoards() {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
        <div className="space-y-1">
          <h1 className="text-3xl font-bold">Funding Rounds</h1>
          <p className="text-muted-foreground">Manage and track your fundraising campaigns</p>
        </div>
        <Dialog>
          <DialogTrigger asChild>
            <Button>
              <Plus className="w-4 h-4 mr-2" />
              New Funding Round
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[425px]">
            <DialogHeader>
              <DialogTitle>Create New Funding Round</DialogTitle>
            </DialogHeader>
            <form className="space-y-6">
              <div className="space-y-2">
                <Label htmlFor="name">Round Name</Label>
                <Input id="name" placeholder="e.g. 2024 Growth Round" />
              </div>

              <div className="space-y-2">
                <Label htmlFor="stage">Funding Stage</Label>
                <Select>
                  <SelectTrigger>
                    <SelectValue placeholder="Select funding stage" />
                  </SelectTrigger>
                  <SelectContent>
                    {FUNDING_STAGES.map((stage) => (
                      <SelectItem key={stage.value} value={stage.value}>
                        <div className="space-y-1">
                          <div className="font-medium">{stage.label}</div>
                          <div className="text-xs text-muted-foreground">{stage.description}</div>
                        </div>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="target">Target Amount ($)</Label>
                <Input id="target" type="number" placeholder="5000000" />
              </div>

              <div className="space-y-2">
                <Label htmlFor="description">Description</Label>
                <Textarea
                  id="description"
                  placeholder="Describe the purpose and goals of this funding round..."
                  rows={4}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="deadline">Deadline</Label>
                <Input id="deadline" type="date" />
              </div>

              <div className="flex gap-4">
                <Button type="submit" className="flex-1">
                  Create Funding Round
                </Button>
                <Button type="button" variant="outline" className="flex-1">
                  Cancel
                </Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 md:gap-6">
        {exampleRounds.map((round) => {
          const statusConfig = getStatusConfig(round.status as RoundStatus)
          const StatusIcon = statusConfig.icon
          const stage = FUNDING_STAGES.find(s => s.value === round.stage)
          
          return (
            <Card key={round.id} className="flex flex-col">
              <CardHeader>
                <div className="flex items-center justify-between mb-2">
                  <div className="space-y-1.5">
                    <CardTitle className="text-xl">{round.name}</CardTitle>
                    <div className="text-sm font-medium text-muted-foreground">
                      {stage?.label}
                    </div>
                  </div>
                  <div className={`flex items-center gap-1.5 px-2.5 py-1.5 rounded-full text-xs font-medium ${statusConfig.textColor} ${statusConfig.bgColor}`}>
                    <StatusIcon className="w-3.5 h-3.5" />
                    {statusConfig.label}
                  </div>
                </div>
                <CardDescription className="line-clamp-2">{round.description}</CardDescription>
              </CardHeader>
              <CardContent className="flex-1">
                <div className="space-y-4">
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span className="text-muted-foreground">Progress</span>
                      <span className="font-medium">
                        {formatCurrency(round.raised)} / {formatCurrency(round.target)}
                      </span>
                    </div>
                    <div className="w-full bg-muted rounded-full h-2">
                      <div 
                        className={`h-2 rounded-full ${statusConfig.color}`}
                        style={{ width: `${(round.raised / round.target) * 100}%` }}
                      />
                    </div>
                  </div>
                  
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span className="text-muted-foreground">Deadline</span>
                      <span>{new Date(round.deadline).toLocaleDateString()}</span>
                    </div>
                    {round.closedReason && (
                      <div className="text-sm text-muted-foreground">
                        <span className="font-medium">Closed: </span>
                        {round.closedReason}
                      </div>
                    )}
                  </div>

                  <Link href={`/fundraising/${round.id}`} className="block mt-4">
                    <Button variant="outline" className="w-full">
                      View Details
                    </Button>
                  </Link>
                </div>
              </CardContent>
            </Card>
          )
        })}
      </div>
    </div>
  )
}
