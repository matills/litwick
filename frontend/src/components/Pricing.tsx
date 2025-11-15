import { Button } from "@/components/ui/button";
import { Check } from "lucide-react";

interface Plan {
  name: string;
  price: string;
  credits: string;
  popular?: boolean;
  priceNote?: string;
  features: string[];
}

const plans: Plan[] = [
  {
    name: "Free",
    price: "0",
    credits: "120",
    features: [
      "120 minutos/mes",
      "Hasta 15 min por archivo",
      "Formatos: TXT, SRT",
      "Retención: 7 días"
    ]
  },
  {
    name: "Pro",
    price: "15",
    credits: "1800",
    popular: true,
    features: [
      "1,800 minutos/mes (30 horas)",
      "Hasta 3 horas por archivo",
      "Formatos: TXT, SRT, VTT, DOCX",
      "Editor con timestamps editables",
      "Retención: 90 días",
      "Soporte por email"
    ]
  },
  {
    name: "Team",
    price: "20",
    credits: "3000",
    priceNote: "/persona (mín. 3)",
    features: [
      "3,000 minutos/persona/mes",
      "Hasta 8 horas por archivo",
      "Workspace compartido",
      "Gestión de equipos y roles",
      "Colaboración en transcripciones",
      "Retención ilimitada",
      "Soporte prioritario"
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
                  <span className="text-muted-foreground">
                    {plan.priceNote || '/mes'}
                  </span>
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
                {plan.name === 'Free' ? 'Comenzar gratis' : 'Comenzar ahora'}
              </Button>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};
