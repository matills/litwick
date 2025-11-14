import { useState } from "react";
import { Upload, FileAudio, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { apiClient } from "@/lib/api";
import { toast } from "sonner";
import { useNavigate } from "react-router-dom";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export const UploadZone = () => {
  const [isDragging, setIsDragging] = useState(false);
  const [file, setFile] = useState<File | null>(null);
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const uploadMutation = useMutation({
    mutationFn: apiClient.uploadFile,
    onSuccess: () => {
      toast.success("Archivo subido exitosamente. Transcripción iniciada.");
      setFile(null);
      queryClient.invalidateQueries({ queryKey: ['dashboard-stats'] });
      queryClient.invalidateQueries({ queryKey: ['transcriptions'] });
      navigate('/dashboard');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || "Error al subir el archivo");
    },
  });

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = () => {
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    const droppedFile = e.dataTransfer.files[0];
    if (droppedFile) setFile(droppedFile);
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0];
    if (selectedFile) setFile(selectedFile);
  };

  return (
    <div className="space-y-4">
      <div
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        className={`relative border-2 border-dashed rounded-2xl p-12 transition-all ${
          isDragging
            ? "border-accent bg-accent/5 scale-[1.02]"
            : "border-border hover:border-accent/50 hover:bg-accent/5"
        }`}
      >
        <div className="flex flex-col items-center justify-center text-center">
          <div className="w-16 h-16 bg-accent/10 rounded-2xl flex items-center justify-center mb-6">
            <Upload className="w-8 h-8 text-accent" />
          </div>
          
          <h3 className="text-xl font-semibold text-foreground mb-2">
            Arrastra tu archivo aquí
          </h3>
          <p className="text-muted-foreground mb-6">
            o haz clic para seleccionar
          </p>

          <input
            type="file"
            id="file-upload"
            className="hidden"
            accept="audio/*,video/*"
            onChange={handleFileSelect}
          />
          <label htmlFor="file-upload">
            <Button asChild className="bg-accent hover:bg-accent/90 text-accent-foreground">
              <span>Seleccionar Archivo</span>
            </Button>
          </label>

          <p className="text-xs text-muted-foreground mt-4">
            Soporta MP3, WAV, MP4, AVI (máx. 500MB)
          </p>
        </div>
      </div>

      {file && (
        <div className="flex items-center gap-4 p-4 bg-card border border-border rounded-xl">
          <div className="w-12 h-12 bg-accent/10 rounded-lg flex items-center justify-center flex-shrink-0">
            <FileAudio className="w-6 h-6 text-accent" />
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium text-foreground truncate">{file.name}</p>
            <p className="text-xs text-muted-foreground">
              {(file.size / 1024 / 1024).toFixed(2)} MB
            </p>
          </div>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setFile(null)}
            className="flex-shrink-0"
          >
            <X className="w-4 h-4" />
          </Button>
        </div>
      )}

      {file && (
        <Button
          className="w-full bg-accent hover:bg-accent/90 text-accent-foreground py-6 text-base font-semibold"
          onClick={() => uploadMutation.mutate(file)}
          disabled={uploadMutation.isPending}
        >
          {uploadMutation.isPending ? "Subiendo..." : "Iniciar Transcripción"}
        </Button>
      )}
    </div>
  );
};
