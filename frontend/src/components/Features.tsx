import { FileAudio, Languages, Download, Clock, Zap, Shield } from "lucide-react";

const features = [
  {
    icon: FileAudio,
    title: "Múltiples Formatos",
    description: "Soporta MP3, WAV, MP4, AVI y más. Sube cualquier archivo de audio o video."
  },
  {
    icon: Clock,
    title: "Transcripción Rápida",
    description: "Procesa tus archivos en minutos con tecnología de IA de última generación."
  },
  {
    icon: Languages,
    title: "Múltiples Idiomas",
    description: "Transcribe audio en diversos idiomas con alta precisión."
  },
  {
    icon: Download,
    title: "Exportación Simple",
    description: "Descarga tus transcripciones en formato .txt y .srt para subtítulos."
  },
  {
    icon: Zap,
    title: "Procesamiento Automático",
    description: "Sube tu archivo y deja que la IA haga el trabajo pesado por ti."
  },
  {
    icon: Shield,
    title: "Seguro y Privado",
    description: "Tus archivos están protegidos y solo tú tienes acceso a tus transcripciones."
  }
];

export const Features = () => {
  return (
    <section className="py-24 bg-background">
      <div className="container px-4 mx-auto">
        <div className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="text-4xl md:text-5xl font-bold text-foreground mb-4">
            Todo lo que necesitas para
            <span className="text-accent"> transcribir mejor</span>
          </h2>
          <p className="text-lg text-muted-foreground">
            Funcionalidades profesionales diseñadas para ahorrarte tiempo y mejorar tu flujo de trabajo.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {features.map((feature, index) => {
            const Icon = feature.icon;
            return (
              <div 
                key={index}
                className="group p-8 bg-card border border-border rounded-2xl hover:shadow-medium hover:border-accent/50 transition-all duration-300 hover:-translate-y-1"
              >
                <div className="w-14 h-14 bg-accent/10 rounded-xl flex items-center justify-center mb-6 group-hover:bg-accent/20 transition-colors">
                  <Icon className="w-7 h-7 text-accent" />
                </div>
                <h3 className="text-xl font-semibold text-foreground mb-3">
                  {feature.title}
                </h3>
                <p className="text-muted-foreground leading-relaxed">
                  {feature.description}
                </p>
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
};
