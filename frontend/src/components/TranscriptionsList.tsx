import { FileText, Download, Eye, Trash2, RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api";
import { toast } from "sonner";
import { format } from "date-fns";
import { es } from "date-fns/locale";
import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { ScrollArea } from "@/components/ui/scroll-area";

export const TranscriptionsList = () => {
  const queryClient = useQueryClient();
  const [previewTranscription, setPreviewTranscription] = useState<any>(null);

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['transcriptions'],
    queryFn: apiClient.getTranscriptions,
    refetchInterval: (query) => {
      const transcriptions = query.state.data?.transcriptions || [];
      const hasProcessing = transcriptions.some((t: any) => t.status === 'processing');
      return hasProcessing ? 5000 : false;
    },
  });

  const deleteMutation = useMutation({
    mutationFn: apiClient.deleteTranscription,
    onSuccess: () => {
      toast.success("Transcripción eliminada");
      queryClient.invalidateQueries({ queryKey: ['transcriptions'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard-stats'] });
    },
    onError: () => {
      toast.error("Error al eliminar la transcripción");
    },
  });

  const handleDownload = async (id: string, format: string) => {
    try {
      const blob = await apiClient.downloadTranscription(id, format);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `transcription-${id}.${format}`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
      toast.success(`Descargando archivo ${format.toUpperCase()}`);
    } catch (error) {
      toast.error("Error al descargar el archivo");
    }
  };

  const transcriptions = data?.transcriptions || [];
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-foreground">Mis Transcripciones</h2>
        <Button variant="outline" size="sm" className="gap-2" onClick={() => refetch()}>
          <RefreshCw className="w-4 h-4" />
          Actualizar
        </Button>
      </div>

      {isLoading && (
        <div className="text-center py-8">
          <p className="text-muted-foreground">Cargando...</p>
        </div>
      )}

      <div className="space-y-3">
        {transcriptions.map((item) => (
          <div
            key={item.id}
            className="p-6 bg-card border border-border rounded-xl hover:shadow-soft hover:border-accent/30 transition-all group"
          >
            <div className="flex items-start gap-4">
              {/* Icon */}
              <div className="w-12 h-12 bg-accent/10 rounded-lg flex items-center justify-center flex-shrink-0">
                <FileText className="w-6 h-6 text-accent" />
              </div>

              {/* Content */}
              <div className="flex-1 min-w-0">
                <div className="flex items-start justify-between gap-4 mb-3">
                  <div className="flex-1 min-w-0">
                    <h3 className="text-base font-semibold text-foreground truncate mb-1">
                      {item.file_name}
                    </h3>
                    <div className="flex flex-wrap items-center gap-3 text-sm text-muted-foreground">
                      <span>Creado: {format(new Date(item.created_at), "dd MMM yyyy, HH:mm", { locale: es })}</span>
                      <span>•</span>
                      <span>Duración: {Math.ceil(item.duration / 60)} min</span>
                      {item.status === 'completed' && item.credits_used > 0 && (
                        <>
                          <span>•</span>
                          <span>Créditos: {item.credits_used} min</span>
                        </>
                      )}
                      {item.language && (
                        <>
                          <span>•</span>
                          <span>Idioma: {item.language.toUpperCase()}</span>
                        </>
                      )}
                    </div>
                  </div>
                  <Badge className={
                    item.status === 'completed'
                      ? "bg-green-500/10 text-green-600 hover:bg-green-500/20 border-green-500/20"
                      : item.status === 'processing'
                      ? "bg-orange-500/10 text-orange-600 hover:bg-orange-500/20 border-orange-500/20"
                      : "bg-red-500/10 text-red-600 hover:bg-red-500/20 border-red-500/20"
                  }>
                    {item.status === 'completed' ? 'Completado' : item.status === 'processing' ? 'Procesando' : 'Error'}
                  </Badge>
                </div>

                {/* Actions */}
                {item.status === 'completed' && (
                  <div className="flex flex-wrap gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      className="gap-2 text-xs"
                      onClick={() => setPreviewTranscription(item)}
                    >
                      <Eye className="w-3.5 h-3.5" />
                      Ver
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      className="gap-2 text-xs"
                      onClick={() => handleDownload(item.id, 'txt')}
                    >
                      <Download className="w-3.5 h-3.5" />
                      TXT
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      className="gap-2 text-xs"
                      onClick={() => handleDownload(item.id, 'srt')}
                    >
                      <Download className="w-3.5 h-3.5" />
                      SRT
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      className="gap-2 text-xs"
                      onClick={() => handleDownload(item.id, 'vtt')}
                    >
                      <Download className="w-3.5 h-3.5" />
                      VTT
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="gap-2 text-xs text-destructive hover:text-destructive hover:bg-destructive/10 ml-auto"
                      onClick={() => deleteMutation.mutate(item.id)}
                      disabled={deleteMutation.isPending}
                    >
                      <Trash2 className="w-3.5 h-3.5" />
                      Eliminar
                    </Button>
                  </div>
                )}
                {item.status === 'processing' && (
                  <div className="text-sm text-muted-foreground italic">
                    Procesando transcripción...
                  </div>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>

      {transcriptions.length === 0 && (
        <div className="py-16 text-center">
          <div className="w-16 h-16 bg-muted rounded-full flex items-center justify-center mx-auto mb-4">
            <FileText className="w-8 h-8 text-muted-foreground" />
          </div>
          <h3 className="text-lg font-semibold text-foreground mb-2">
            No hay transcripciones aún
          </h3>
          <p className="text-muted-foreground">
            Sube tu primer archivo de audio o video para comenzar
          </p>
        </div>
      )}

      {/* Preview Modal */}
      <Dialog open={!!previewTranscription} onOpenChange={() => setPreviewTranscription(null)}>
        <DialogContent className="max-w-4xl max-h-[80vh]">
          <DialogHeader>
            <DialogTitle>{previewTranscription?.file_name}</DialogTitle>
            <DialogDescription>
              Transcripción completada el {previewTranscription?.completed_at && format(new Date(previewTranscription.completed_at), "dd MMM yyyy, HH:mm", { locale: es })}
            </DialogDescription>
          </DialogHeader>
          <ScrollArea className="max-h-[60vh] pr-4">
            <div className="space-y-4">
              <div className="text-sm text-foreground whitespace-pre-wrap">
                {previewTranscription?.transcript_text || "No hay texto disponible"}
              </div>
            </div>
          </ScrollArea>
          <div className="flex gap-2 justify-end">
            <Button
              variant="outline"
              onClick={() => handleDownload(previewTranscription?.id, 'txt')}
            >
              <Download className="w-4 h-4 mr-2" />
              Descargar TXT
            </Button>
            <Button
              variant="outline"
              onClick={() => handleDownload(previewTranscription?.id, 'srt')}
            >
              <Download className="w-4 h-4 mr-2" />
              Descargar SRT
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
};
