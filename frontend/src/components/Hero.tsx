import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Upload, Sparkles, Zap } from "lucide-react";
import { AuthModal } from "./AuthModal";

export const Hero = () => {
  const [authModalOpen, setAuthModalOpen] = useState(false);

  return (
    <>
      <AuthModal open={authModalOpen} onOpenChange={setAuthModalOpen} />
    <section className="relative min-h-screen flex items-center justify-center overflow-hidden bg-gradient-hero">
      {/* Animated background elements */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-accent/10 rounded-full blur-3xl animate-pulse" />
        <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-primary/10 rounded-full blur-3xl animate-pulse delay-1000" />
      </div>

      <div className="container relative z-10 px-4 mx-auto">
        <div className="max-w-4xl mx-auto text-center">
          <div className="inline-flex items-center gap-2 px-4 py-2 mb-8 bg-accent/10 backdrop-blur-sm border border-accent/20 rounded-full">
            <Sparkles className="w-4 h-4 text-accent" />
            <span className="text-sm font-medium text-accent">Transcripciones con IA en segundos</span>
          </div>

          <h1 className="mb-6 text-5xl md:text-7xl font-bold text-white leading-tight">
            Convierte Audio en
            <span className="block bg-gradient-to-r from-accent to-accent/70 bg-clip-text text-transparent">
              Texto Perfecto
            </span>
          </h1>

          <p className="mb-10 text-xl md:text-2xl text-white/80 max-w-2xl mx-auto">
            Transcripciones automáticas precisas con timestamps, edición en tiempo real y exportación instantánea.
          </p>

          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
            <Button
              size="lg"
              className="bg-accent hover:bg-accent/90 text-accent-foreground shadow-accent text-lg px-8 py-6 rounded-xl font-semibold transition-all hover:scale-105"
              onClick={() => setAuthModalOpen(true)}
            >
              <Upload className="mr-2 h-5 w-5" />
              Comenzar Gratis
            </Button>
            <Button
              size="lg"
              variant="outline"
              className="bg-white/10 backdrop-blur-sm border-white/20 text-white hover:bg-white/20 text-lg px-8 py-6 rounded-xl font-semibold transition-all"
              onClick={() => setAuthModalOpen(true)}
            >
              Ver Demo
              <Zap className="ml-2 h-5 w-5" />
            </Button>
          </div>

          <div className="mt-16 grid grid-cols-1 md:grid-cols-3 gap-8 max-w-3xl mx-auto">
            {[
              { value: "98%", label: "Precisión" },
              { value: "50+", label: "Idiomas" },
              { value: "10x", label: "Más Rápido" }
            ].map((stat, index) => (
              <div key={index} className="p-6 bg-white/5 backdrop-blur-sm border border-white/10 rounded-2xl">
                <div className="text-4xl font-bold text-accent mb-2">{stat.value}</div>
                <div className="text-white/70">{stat.label}</div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
    </>
  );
};
