import { DashboardLayout } from "@/components/DashboardLayout";
import { UploadZone } from "@/components/UploadZone";

const Upload = () => {
  return (
    <DashboardLayout>
      <div className="max-w-3xl mx-auto">
        <h1 className="text-3xl font-bold text-foreground mb-2">Subir Archivo</h1>
        <p className="text-muted-foreground mb-8">
          Sube tu archivo de audio o video para comenzar la transcripción
        </p>
        <UploadZone />
      </div>
    </DashboardLayout>
  );
};

export default Upload;
