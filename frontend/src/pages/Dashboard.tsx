import { DashboardLayout } from "@/components/DashboardLayout";
import { UploadZone } from "@/components/UploadZone";
import { StatsCards } from "@/components/StatsCards";
import { TranscriptionsList } from "@/components/TranscriptionsList";

const Dashboard = () => {
  return (
    <DashboardLayout>
      <div className="max-w-7xl mx-auto space-y-8">
        {/* Stats */}
        <StatsCards />

        {/* Upload Zone */}
        <div>
          <h2 className="text-2xl font-bold text-foreground mb-4">Nueva Transcripción</h2>
          <UploadZone />
        </div>

        {/* Transcriptions List */}
        <TranscriptionsList />
      </div>
    </DashboardLayout>
  );
};

export default Dashboard;
