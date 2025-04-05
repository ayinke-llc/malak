"use client"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Textarea } from "@/components/ui/textarea"
import { CheckCircle2, Clock, Plus, XCircle } from "lucide-react"
import Link from "next/link"
import * as yup from "yup"
import { useForm } from "react-hook-form"
import { yupResolver } from "@hookform/resolvers/yup"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { CREATE_FUNDRAISING_PIPELINE } from "@/lib/query-constants"
import client from "@/lib/client"
import { MalakFundraisePipelineStage, ServerAPIStatus } from "@/client/Api"
import type { ServerCreateNewPipelineRequest } from "@/client/Api"
import { toast } from "sonner"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form"
import React from "react"
import { NumericFormat } from "react-number-format"
import { AxiosError } from "axios"

type FundingStage = MalakFundraisePipelineStage;

const FUNDING_STAGES: { value: MalakFundraisePipelineStage; label: string; description: string }[] = [
  {
    value: MalakFundraisePipelineStage.FundraisePipelineStageFamilyAndFriend,
    label: "Family & Friends",
    description: "Initial funding from close connections"
  },
  {
    value: MalakFundraisePipelineStage.FundraisePipelineStagePreSeed,
    label: "Pre-Seed",
    description: "Very early stage funding to develop initial product"
  },
  {
    value: MalakFundraisePipelineStage.FundraisePipelineStageSeed,
    label: "Seed",
    description: "Early stage funding to validate product market fit"
  },
  {
    value: MalakFundraisePipelineStage.FundraisePipelineStageSeriesA,
    label: "Series A",
    description: "First significant round of venture capital financing"
  },
  {
    value: MalakFundraisePipelineStage.FundraisePipelineStageSeriesB,
    label: "Series B",
    description: "Funding for business development and market growth"
  },
  {
    value: MalakFundraisePipelineStage.FundraisePipelineStageSeriesC,
    label: "Series C",
    description: "Scaling, expansion, and possible acquisitions"
  },
  {
    value: MalakFundraisePipelineStage.FundraisePipelineStageBridgeRound,
    label: "Bridge",
    description: "Short-term funding between major rounds"
  }
];

// Example data - replace with actual data from your API
const exampleRounds = [
  {
    id: 1,
    name: "2024 Growth Round",
    stage: MalakFundraisePipelineStage.FundraisePipelineStageSeriesA,
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
    stage: MalakFundraisePipelineStage.FundraisePipelineStageSeed,
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
    stage: MalakFundraisePipelineStage.FundraisePipelineStageBridgeRound,
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
    stage: MalakFundraisePipelineStage.FundraisePipelineStageSeriesB,
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
    stage: MalakFundraisePipelineStage.FundraisePipelineStageSeriesC,
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
    stage: MalakFundraisePipelineStage.FundraisePipelineStagePreSeed,
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

type FormData = {
  title: string;
  stage: MalakFundraisePipelineStage;
  amount: number;
  description: string;
  expected_close_date: string;
  start_date: string;
};

const createPipelineSchema = yup.object().shape({
  title: yup.string().required("Round name is required"),
  stage: yup.mixed<MalakFundraisePipelineStage>()
    .required("Funding stage is required")
    .oneOf(
      Object.values(MalakFundraisePipelineStage),
      "Invalid funding stage"
    ),
  amount: yup.number().required("Target amount is required").positive("Target must be positive"),
  description: yup.string().required("Description is required"),
  expected_close_date: yup.string().required("Deadline is required"),
  start_date: yup.string().required("Start date is required")
    .test("start-before-end", "Start date must be before deadline", function (startDate) {
      const deadline = this.parent.expected_close_date;
      if (!startDate || !deadline) return true;
      return new Date(startDate) < new Date(deadline);
    })
}) satisfies yup.ObjectSchema<FormData>;

const NumericInput = React.forwardRef<HTMLInputElement, any>((props, ref) => {
  const { onChange, ...other } = props;

  return (
    <NumericFormat
      {...other}
      getInputRef={ref}
      thousandSeparator={true}
      prefix="$"
      onValueChange={(values) => {
        onChange(values.floatValue || 0);
      }}
      valueIsNumericString
      className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
    />
  );
});
NumericInput.displayName = "NumericInput";

// Add this helper function before the FundraisingBoards component
const getTomorrowDate = () => {
  const tomorrow = new Date();
  tomorrow.setDate(tomorrow.getDate() + 1);
  return tomorrow.toISOString().split('T')[0];
};

const getThreeMonthsFromNow = () => {
  const threeMonths = new Date();
  threeMonths.setMonth(threeMonths.getMonth() + 3);
  return threeMonths.toISOString().split('T')[0];
};

export default function FundraisingBoards() {
  const queryClient = useQueryClient();
  const [open, setOpen] = React.useState(false);
  
  const form = useForm<FormData>({
    resolver: yupResolver(createPipelineSchema),
    defaultValues: {
      title: "",
      stage: undefined,
      amount: undefined,
      description: "",
      expected_close_date: getThreeMonthsFromNow(),
      start_date: getTomorrowDate()
    }
  });

  const createPipeline = useMutation({
    mutationKey: [CREATE_FUNDRAISING_PIPELINE],
    mutationFn: async (data: FormData): Promise<unknown> => {
      const request: ServerCreateNewPipelineRequest = {
        ...data,
        amount: data.amount * 100,
        start_date: Math.floor(new Date(data.start_date).getTime() / 1000),
        expected_close_date: Math.floor(new Date(data.expected_close_date).getTime() / 1000)
      };
      return client.pipelines.pipelinesCreate(request);
    },
    onSuccess: () => {
      toast.success("Funding round created successfully");
      setOpen(false);
      form.reset({
        title: "",
        stage: undefined,
        amount: undefined,
        description: "",
        expected_close_date: getThreeMonthsFromNow(),
        start_date: getTomorrowDate()
      });
      queryClient.invalidateQueries({ queryKey: ["fundraisingPipelines"] });
    },
    onError: (err: AxiosError<ServerAPIStatus>) => {
      toast.error(err.response?.data.message ?? "An error occurred while creating funding round");
      console.error("Error creating funding round:", err);
    }
  });

  const onSubmit = (data: FormData) => {
    createPipeline.mutate(data);
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
        <div className="space-y-1">
          <h1 className="text-3xl font-bold">Funding Rounds</h1>
          <p className="text-muted-foreground">Manage and track your fundraising campaigns</p>
        </div>
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="w-4 h-4 mr-2" />
              New Funding Round
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[800px]">
            <DialogHeader>
              <DialogTitle>Create New Funding Round</DialogTitle>
            </DialogHeader>
            <Form {...form}>
              <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
                <div className="grid grid-cols-2 gap-6">
                  <FormField
                    control={form.control}
                    name="title"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Round Name</FormLabel>
                        <FormControl>
                          <Input placeholder="e.g. 2024 Growth Round" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="stage"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Funding Stage</FormLabel>
                        <Select onValueChange={field.onChange} value={field.value}>
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder="Select funding stage" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {FUNDING_STAGES.map((stage) => (
                              <SelectItem 
                                key={stage.value} 
                                value={stage.value}
                                className="py-2 cursor-pointer focus:bg-accent/30 hover:bg-accent/30"
                              >
                                <div className="space-y-0.5">
                                  <div className="font-medium text-foreground">{stage.label}</div>
                                  <div className="text-xs text-muted-foreground/80">{stage.description}</div>
                                </div>
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="amount"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Target Amount ($)</FormLabel>
                        <FormControl>
                          <NumericInput
                            placeholder="Enter amount"
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <div className="grid grid-cols-2 gap-4">
                    <FormField
                      control={form.control}
                      name="start_date"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Start Date</FormLabel>
                          <FormControl>
                            <Input 
                              type="date" 
                              min={getTomorrowDate()}
                              {...field} 
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="expected_close_date"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Expected Close Date</FormLabel>
                          <FormControl>
                            <Input 
                              type="date" 
                              min={getTomorrowDate()}
                              {...field} 
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                </div>

                <FormField
                  control={form.control}
                  name="description"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Description</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder="Describe the purpose and goals of this funding round..."
                          className="min-h-[100px]"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <div className="flex gap-4 justify-end">
                  <Button type="button" variant="outline" onClick={() => setOpen(false)}>
                    Cancel
                  </Button>
                  <Button type="submit" loading={createPipeline.isPending}>
                    {createPipeline.isPending ? "Creating..." : "Create Funding Round"}
                  </Button>
                </div>
              </form>
            </Form>
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
