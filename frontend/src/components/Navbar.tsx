import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Flame } from "lucide-react";
import { AuthModal } from "./AuthModal";

export const Navbar = () => {
  const [authModalOpen, setAuthModalOpen] = useState(false);

  return (
    <>
      <nav className="fixed top-0 left-0 right-0 z-50 bg-background/80 backdrop-blur-md border-b border-border">
        <div className="container px-4 mx-auto">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center gap-2">
              <div className="w-10 h-10 bg-gradient-to-br from-orange-400 to-orange-600 rounded-lg flex items-center justify-center shadow-lg">
                <Flame className="w-6 h-6 text-white" />
              </div>
              <span className="text-xl font-bold text-foreground">Litwick</span>
            </div>

            <div className="hidden md:flex items-center gap-8">
              <a href="#features" className="text-foreground hover:text-accent transition-colors font-medium">
                Funcionalidades
              </a>
              <a href="#pricing" className="text-foreground hover:text-accent transition-colors font-medium">
                Precios
              </a>
              <a href="#docs" className="text-foreground hover:text-accent transition-colors font-medium">
                Docs
              </a>
            </div>

            <div className="flex items-center gap-4">
              <Button
                variant="ghost"
                className="text-foreground hover:text-accent"
                onClick={() => setAuthModalOpen(true)}
              >
                Iniciar Sesión
              </Button>
              <Button
                className="bg-accent hover:bg-accent/90 text-accent-foreground shadow-accent"
                onClick={() => setAuthModalOpen(true)}
              >
                Comenzar Gratis
              </Button>
            </div>
          </div>
        </div>
      </nav>

      <AuthModal open={authModalOpen} onOpenChange={setAuthModalOpen} />
    </>
  );
};
