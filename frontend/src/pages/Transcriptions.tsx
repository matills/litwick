import { DashboardLayout } from "@/components/DashboardLayout";
import { TranscriptionsList } from "@/components/TranscriptionsList";

const Transcriptions = () => {
  return (
    <DashboardLayout>
      <div className="max-w-7xl mx-auto">
        <TranscriptionsList />
      </div>
    </DashboardLayout>
  );
};

export default Transcriptions;
