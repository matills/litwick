import { Button } from "@/components/ui/button";
import { Check } from "lucide-react";
import { useState } from "react";
import { AuthModal } from "./AuthModal";

const plans = [
  {
    name: "Básico",
    price: "5",
    credits: "120",
    discount: 0,
    features: [
      "120 minutos de transcripción",
      "Múltiples idiomas",
      "Exportar en .txt y .srt",
      "Soporte por email"
    ]
  },
  {
    name: "Estándar",
    price: "10",
    credits: "300",
    discount: 17,
    popular: true,
    features: [
      "300 minutos de transcripción",
      "Múltiples idiomas",
      "Exportar en .txt y .srt",
      "Soporte por email",
      "Ahorrás 17%"
    ]
  },
  {
    name: "Premium",
    price: "18",
    credits: "600",
    discount: 25,
    features: [
      "600 minutos de transcripción",
      "Múltiples idiomas",
      "Exportar en .txt y .srt",
      "Soporte por email",
      "Ahorrás 25%"
    ]
  },
  {
    name: "Max",
    price: "40",
    credits: "1500",
    discount: 33,
    features: [
      "1500 minutos de transcripción",
      "Múltiples idiomas",
      "Exportar en .txt y .srt",
      "Soporte prioritario",
      "Ahorrás 33%"
    ]
  }
];

export const Pricing = () => {
  const [authModalOpen, setAuthModalOpen] = useState(false);

  return (
    <>
      <section className="py-24 bg-secondary/30">
        <div className="container px-4 mx-auto">
          <div className="text-center max-w-3xl mx-auto mb-16">
            <h2 className="text-4xl md:text-5xl font-bold text-foreground mb-4">
              Paquetes de créditos
            </h2>
            <p className="text-lg text-muted-foreground">
              Compra minutos cuando los necesites. Sin suscripciones, sin compromisos.
            </p>
          </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 max-w-7xl mx-auto">
          {plans.map((plan, index) => (
            <div 
              key={index}
              className={`relative p-8 bg-card border rounded-2xl transition-all duration-300 hover:shadow-medium ${
                plan.popular 
                  ? 'border-accent shadow-accent scale-105 md:scale-110' 
                  : 'border-border hover:border-accent/50'
              }`}
            >
              {plan.popular && (
                <div className="absolute -top-4 left-1/2 -translate-x-1/2 px-4 py-1 bg-gradient-accent text-accent-foreground text-sm font-semibold rounded-full">
                  Más Popular
                </div>
              )}

              <div className="mb-6">
                <h3 className="text-2xl font-bold text-foreground mb-2">{plan.name}</h3>
                <div className="flex items-baseline gap-2">
                  <span className="text-5xl font-bold text-foreground">${plan.price}</span>
                </div>
                {plan.discount > 0 && (
                  <p className="text-sm text-accent font-medium mt-1">Ahorrás {plan.discount}%</p>
                )}
                <p className="text-sm text-muted-foreground mt-2">{plan.credits} minutos</p>
              </div>

              <ul className="space-y-4 mb-8">
                {plan.features.map((feature, fIndex) => (
                  <li key={fIndex} className="flex items-start gap-3">
                    <div className="w-5 h-5 rounded-full bg-accent/10 flex items-center justify-center flex-shrink-0 mt-0.5">
                      <Check className="w-3 h-3 text-accent" />
                    </div>
                    <span className="text-card-foreground text-sm">{feature}</span>
                  </li>
                ))}
              </ul>

              <Button
                className={`w-full ${
                  plan.popular
                    ? 'bg-accent hover:bg-accent/90 text-accent-foreground shadow-accent'
                    : 'bg-primary hover:bg-primary/90'
                }`}
                onClick={() => setAuthModalOpen(true)}
              >
                Comenzar ahora
              </Button>
            </div>
          ))}
        </div>
      </div>
    </section>

    <AuthModal open={authModalOpen} onOpenChange={setAuthModalOpen} />
    </>
  );
};
