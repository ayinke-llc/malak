"use client";

import ListDashboards from "@/components/ui/dashboards/list";
import CreateDashboardModal from "@/components/ui/dashboards/create-modal";

export default function Page() {
  return (
    <>
      <div className="pt-6 bg-background">
        <section>
          <div className="sm:flex sm:items-center sm:justify-between">
            <div>
              <h3
                id="company-dashboards"
                className="text-lg font-medium"
              >
                Your dashboards
              </h3>
              <p className="text-sm text-muted-foreground">
                View and manage dashboards created from the data of your integrations
              </p>
            </div>

            <div>
              <CreateDashboardModal />
            </div>
          </div>
        </section>

        <section className="mt-10">
          <ListDashboards />
        </section>
      </div>
    </>
  );
}
