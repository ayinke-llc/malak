import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

export default function BannedPage() {
  return (
    <div className="fixed inset-0 bg-background flex items-center justify-center">
      <div className="absolute inset-0 bg-red-50/50" />
      <Card className="relative w-full max-w-lg mx-4 border-2 border-red-200 shadow-lg">
        <CardHeader className="border-b border-red-100">
          <CardTitle className="text-4xl text-center text-red-600 font-bold">Access Denied</CardTitle>
        </CardHeader>
        <CardContent className="space-y-6 py-8">
          <div className="flex justify-center mb-4">
            <svg className="w-16 h-16 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
          <p className="text-2xl font-semibold text-center text-gray-900">
            Your account has been banned from accessing the dashboard.
          </p>
          <p className="text-lg text-center text-gray-600">
            If you believe this is a mistake, please contact support@ayinke.ventures
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
