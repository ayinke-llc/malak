"use client"

import * as React from "react"
import * as TabsPrimitive from "@radix-ui/react-tabs"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "@/lib/utils"

const Tabs = TabsPrimitive.Root

const TabsListVariants = cva(
  "inline-flex items-center justify-start h-9",
  {
    variants: {
      variant: {
        default: "border-b border-border bg-background",
        vercel: "border-b border-border bg-background",
      },
      size: {
        default: "h-9",
        sm: "h-8 text-xs",
        lg: "h-10",
        icon: "h-9 w-9",
      },
      width: {
        default: "w-full",
        fit: "w-fit"
      }
    },
    defaultVariants: {
      variant: "default",
      size: "default",
      width: "default"
    },
  }
)

const TabsTriggerVariants = cva(
  "inline-flex items-center justify-center whitespace-nowrap text-sm font-medium transition-all disabled:pointer-events-none border-b-2 border-transparent",
  {
    variants: {
      variant: {
        default: "px-4 py-2 h-9 data-[state=active]:border-primary text-muted-foreground data-[state=active]:text-foreground hover:text-foreground",
        vercel: "px-4 py-2 h-9 -mb-px data-[state=active]:border-primary text-muted-foreground data-[state=active]:text-foreground hover:text-foreground",
      },
      size: {
        default: "",
        sm: "text-xs h-8",
        lg: "text-base h-10",
        icon: "h-9 w-9",
      },
      width: {
        default: "w-full",
        fit: "w-fit"
      }
    },
    defaultVariants: {
      variant: "default",
      size: "default",
      width: "default"
    },
  }
)

export interface TabsListProps
  extends React.ComponentPropsWithoutRef<typeof TabsPrimitive.List>,
    VariantProps<typeof TabsListVariants> {
  asChild?: boolean
}

const TabsList = React.forwardRef<
  React.ElementRef<typeof TabsPrimitive.List>,
  TabsListProps
>(({ className, variant, size, width, ...props }, ref) => (
  <TabsPrimitive.List
    ref={ref}
    className={cn(TabsListVariants({ variant, size, width, className }))}
    {...props}
  />
))
TabsList.displayName = TabsPrimitive.List.displayName

export interface TabsTriggerProps
  extends React.ComponentPropsWithoutRef<typeof TabsPrimitive.Trigger>,
    VariantProps<typeof TabsTriggerVariants> {
  asChild?: boolean
}

const TabsTrigger = React.forwardRef<
  React.ElementRef<typeof TabsPrimitive.Trigger>,
  TabsTriggerProps
>(({ className, variant, size, width, ...props }, ref) => (
  <TabsPrimitive.Trigger
    ref={ref}
    className={cn(TabsTriggerVariants({ variant, size, width, className }))}
    {...props}
  />
))
TabsTrigger.displayName = TabsPrimitive.Trigger.displayName

const TabsContent = React.forwardRef<
  React.ElementRef<typeof TabsPrimitive.Content>,
  React.ComponentPropsWithoutRef<typeof TabsPrimitive.Content>
>(({ className, ...props }, ref) => (
  <TabsPrimitive.Content
    ref={ref}
    className={cn(
      "mt-2 ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
      className
    )}
    {...props}
  />
))
TabsContent.displayName = TabsPrimitive.Content.displayName

export { Tabs, TabsList, TabsTrigger, TabsContent }
