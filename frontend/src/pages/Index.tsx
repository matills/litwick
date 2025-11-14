import { Navbar } from "@/components/Navbar";
import { Hero } from "@/components/Hero";
import { Features } from "@/components/Features";
import { Pricing } from "@/components/Pricing";

const Index = () => {
  return (
    <div className="min-h-screen">
      <Navbar />
      <Hero />
      <div id="features">
        <Features />
      </div>
      <div id="pricing">
        <Pricing />
      </div>
      
      {/* Footer */}
      <footer className="py-12 bg-primary text-primary-foreground">
        <div className="container px-4 mx-auto text-center">
          <p className="text-sm opacity-80">
            © 2025 Litwick. Convirtiendo audio en texto con inteligencia artificial.
          </p>
        </div>
      </footer>
    </div>
  );
};

export default Index;
