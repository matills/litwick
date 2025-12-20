import { useEffect, useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Loader2, ExternalLink, XCircle } from "lucide-react";
import { Button } from "@/components/ui/button";

interface PaymentModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  preferenceId: string | null;
  initPoint?: string | null;
  packageName?: string;
  packagePrice?: number;
  packageCredits?: number;
  packageId?: string;
}

export const PaymentModal = ({
  open,
  onOpenChange,
  preferenceId,
  initPoint,
  packageName,
  packagePrice,
  packageCredits,
  packageId,
}: PaymentModalProps) => {
  const [paymentOpened, setPaymentOpened] = useState(false);

  useEffect(() => {
    if (open && initPoint && !paymentOpened) {
      // Open MercadoPago in a new tab
      window.open(initPoint, '_blank', 'noopener,noreferrer');
      setPaymentOpened(true);
    }
  }, [open, initPoint, paymentOpened]);

  useEffect(() => {
    if (!open) {
      setPaymentOpened(false);
    }
  }, [open]);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Procesando pago</DialogTitle>
          <DialogDescription>
            Se ha abierto MercadoPago en una nueva pestaña
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6 py-4">
          {packageName && packagePrice && packageCredits && (
            <div className="bg-secondary/30 rounded-lg p-4 space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Paquete:</span>
                <span className="font-semibold text-foreground">{packageName}</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Créditos:</span>
                <span className="font-semibold text-foreground">{packageCredits} minutos</span>
              </div>
              <div className="h-px bg-border my-2" />
              <div className="flex justify-between items-center">
                <span className="text-base font-medium text-foreground">Total:</span>
                <span className="font-bold text-2xl text-accent">
                  ${packagePrice.toLocaleString('es-AR')} USD
                </span>
              </div>
            </div>
          )}

          <div className="bg-blue-50 dark:bg-blue-950/20 border border-blue-200 dark:border-blue-900 rounded-lg p-4">
            <div className="flex items-start gap-3">
              <ExternalLink className="w-5 h-5 text-blue-500 flex-shrink-0 mt-0.5" />
              <div className="space-y-2">
                <p className="text-sm font-medium text-blue-900 dark:text-blue-100">
                  Completa tu pago en la nueva pestaña
                </p>
                <p className="text-xs text-blue-700 dark:text-blue-300">
                  Después de completar el pago, vuelve aquí. Tus créditos se agregarán automáticamente.
                </p>
              </div>
            </div>
          </div>

          {initPoint && (
            <Button
              variant="outline"
              className="w-full"
              onClick={() => window.open(initPoint, '_blank', 'noopener,noreferrer')}
            >
              <ExternalLink className="w-4 h-4 mr-2" />
              Abrir pago nuevamente
            </Button>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
};
