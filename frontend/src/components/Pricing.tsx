import { Button } from "@/components/ui/button";
import { Check } from "lucide-react";

const plans = [
  {
    name: "Starter",
    price: "9",
    credits: "300",
    features: [
      "300 minutos de transcripción",
      "Hasta 10 idiomas",
      "Exportar en .txt y .srt",
      "Soporte por email"
    ]
  },
  {
    name: "Pro",
    price: "29",
    credits: "1200",
    popular: true,
    features: [
      "1200 minutos de transcripción",
      "50+ idiomas",
      "Todos los formatos de exportación",
      "Editor avanzado con IA",
      "Múltiples speakers",
      "Soporte prioritario"
    ]
  },
  {
    name: "Enterprise",
    price: "99",
    credits: "5000",
    features: [
      "5000 minutos de transcripción",
      "API access",
      "Integración con YouTube/Vimeo",
      "Plantillas personalizadas",
      "Manager dedicado",
      "SLA garantizado"
    ]
  }
];

export const Pricing = () => {
  return (
    <section className="py-24 bg-secondary/30">
      <div className="container px-4 mx-auto">
        <div className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="text-4xl md:text-5xl font-bold text-foreground mb-4">
            Planes que se ajustan a ti
          </h2>
          <p className="text-lg text-muted-foreground">
            Comienza gratis y escala según tus necesidades. Sin compromisos.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-6xl mx-auto">
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
                  <span className="text-muted-foreground">/mes</span>
                </div>
                <p className="text-sm text-muted-foreground mt-2">{plan.credits} minutos incluidos</p>
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
              >
                Comenzar ahora
              </Button>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};
