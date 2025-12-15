import { DashboardLayout } from "@/components/DashboardLayout";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Flame, TrendingUp, Clock, CheckCircle2, Sparkles } from "lucide-react";
import { Separator } from "@/components/ui/separator";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api";
import { format } from "date-fns";
import { es } from "date-fns/locale";
import { toast } from "sonner";
import { useEffect } from "react";

const Credits = () => {
  const queryClient = useQueryClient();

  const { data: dashboardData, isLoading } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: apiClient.getDashboard,
  });

  const { data: transcriptionsData } = useQuery({
    queryKey: ['transcriptions'],
    queryFn: apiClient.getTranscriptions,
  });

  const { data: packagesData } = useQuery({
    queryKey: ['credit-packages'],
    queryFn: apiClient.getCreditPackages,
  });

  const { data: paymentsData } = useQuery({
    queryKey: ['payment-history'],
    queryFn: apiClient.getPaymentHistory,
  });

  const stats = dashboardData?.stats;
  const transcriptions = transcriptionsData?.transcriptions || [];
  const packages = packagesData?.packages || [];
  const payments = paymentsData?.payments || [];

  // Calculate real stats from transcriptions
  const completedTranscriptions = transcriptions.filter((t: any) => t.status === 'completed');
  const totalMinutesUsed = stats?.total_minutes_used || 0;
  const remainingCredits = stats?.credits_remaining || 0;
  const totalCredits = 300; // Free tier default
  const usagePercentage = totalCredits > 0 ? ((totalCredits - remainingCredits) / totalCredits) * 100 : 0;

  // Average minutes per transcription
  const avgMinutes = completedTranscriptions.length > 0
    ? Math.round(totalMinutesUsed / completedTranscriptions.length)
    : 0;

  // Check for payment status in URL
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const status = params.get('status');
    const paymentId = params.get('payment_id');

    if (status && paymentId) {
      // Process payment success/failure
      apiClient.processPaymentSuccess(params)
        .then((data) => {
          if (status === 'success' || status === 'approved') {
            toast.success(`¡Pago exitoso! Se agregaron ${data.payment.credits_amount} créditos a tu cuenta`);
            queryClient.invalidateQueries({ queryKey: ['dashboard-stats'] });
            queryClient.invalidateQueries({ queryKey: ['payment-history'] });
          } else if (status === 'failure' || status === 'rejected') {
            toast.error('El pago fue rechazado. Por favor intenta nuevamente.');
          } else if (status === 'pending') {
            toast.info('Tu pago está pendiente de aprobación.');
          }
          // Clean URL
          window.history.replaceState({}, '', '/credits');
        })
        .catch((err) => {
          console.error('Error processing payment:', err);
          toast.error('Error al procesar el pago');
        });
    }
  }, [queryClient]);

  // Mutation to create payment
  const createPaymentMutation = useMutation({
    mutationFn: apiClient.createPayment,
    onSuccess: (data) => {
      // Redirect to MercadoPago checkout
      window.location.href = data.init_point;
    },
    onError: () => {
      toast.error("Error al crear el pago. Por favor intenta nuevamente.");
    },
  });

  const handleBuyPackage = (packageId: string) => {
    createPaymentMutation.mutate(packageId);
  };

  if (isLoading) {
    return (
      <DashboardLayout>
        <div className="max-w-6xl mx-auto space-y-8">
          <p className="text-center text-muted-foreground">Cargando...</p>
        </div>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <div className="max-w-6xl mx-auto space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-foreground mb-2">Créditos</h1>
          <p className="text-muted-foreground">
            Administra tus minutos de transcripción disponibles
          </p>
        </div>

        {/* Current Balance */}
        <div className="bg-gradient-accent p-8 rounded-2xl text-accent-foreground">
          <div className="flex items-start justify-between mb-6">
            <div>
              <p className="text-sm opacity-90 mb-2">Créditos disponibles</p>
              <div className="flex items-baseline gap-2">
                <span className="text-5xl font-bold">{remainingCredits}</span>
                <span className="text-2xl opacity-75">minutos</span>
              </div>
              <p className="text-sm opacity-75 mt-2">de {totalCredits} minutos totales</p>
            </div>
            <div className="w-16 h-16 bg-white/20 backdrop-blur-sm rounded-xl flex items-center justify-center">
              <Flame className="w-8 h-8" />
            </div>
          </div>

          <div className="space-y-2">
            <Progress value={100 - usagePercentage} className="h-2 bg-white/20" />
            <div className="flex justify-between text-sm opacity-75">
              <span>{totalMinutesUsed} min usados</span>
              <span>{remainingCredits} min restantes</span>
            </div>
          </div>
        </div>

        {/* Quick Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Card className="p-6 border-border">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-accent/10 rounded-lg flex items-center justify-center">
                <TrendingUp className="w-6 h-6 text-accent" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Total usado</p>
                <p className="text-2xl font-bold text-foreground">{totalMinutesUsed} min</p>
              </div>
            </div>
          </Card>

          <Card className="p-6 border-border">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-green-500/10 rounded-lg flex items-center justify-center">
                <CheckCircle2 className="w-6 h-6 text-green-500" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Completadas</p>
                <p className="text-2xl font-bold text-foreground">{stats?.completed_count || 0}</p>
              </div>
            </div>
          </Card>

          <Card className="p-6 border-border">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-orange-500/10 rounded-lg flex items-center justify-center">
                <Clock className="w-6 h-6 text-orange-500" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Promedio</p>
                <p className="text-2xl font-bold text-foreground">{avgMinutes} min</p>
              </div>
            </div>
          </Card>
        </div>

        {/* Buy Credits */}
        <div className="bg-card border border-border rounded-xl p-6 space-y-6">
          <div>
            <h2 className="text-2xl font-bold text-foreground mb-2">Comprar más créditos</h2>
            <p className="text-muted-foreground">Elige el pack que mejor se ajuste a tus necesidades</p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            {packages.map((pkg: any) => (
              <div
                key={pkg.id}
                className={`relative p-6 border rounded-xl transition-all hover:shadow-medium ${
                  pkg.popular
                    ? "border-accent bg-accent/5"
                    : "border-border hover:border-accent/30"
                }`}
              >
                {pkg.popular && (
                  <div className="absolute -top-3 left-1/2 -translate-x-1/2 px-3 py-1 bg-accent text-accent-foreground text-xs font-semibold rounded-full flex items-center gap-1">
                    <Sparkles className="w-3 h-3" />
                    Más popular
                  </div>
                )}

                <div className="text-center mb-4">
                  <h3 className="text-lg font-semibold text-foreground mb-1">{pkg.name}</h3>
                  <p className="text-xs text-muted-foreground mb-3">{pkg.description}</p>
                  <div className="flex items-baseline justify-center gap-1 mb-1">
                    <span className="text-3xl font-bold text-foreground">
                      ${pkg.price.toLocaleString('es-AR')}
                    </span>
                  </div>
                  {pkg.discount > 0 && (
                    <span className="text-xs text-accent font-medium">
                      Ahorrás {pkg.discount}%
                    </span>
                  )}
                </div>

                <div className="mb-6 text-center py-3 bg-secondary/50 rounded-lg">
                  <p className="text-2xl font-bold text-foreground">{pkg.credits}</p>
                  <p className="text-sm text-muted-foreground">minutos</p>
                </div>

                <Button
                  className={`w-full ${
                    pkg.popular
                      ? "bg-accent hover:bg-accent/90 text-accent-foreground"
                      : "bg-primary hover:bg-primary/90"
                  }`}
                  onClick={() => handleBuyPackage(pkg.id)}
                  disabled={createPaymentMutation.isPending}
                >
                  {createPaymentMutation.isPending ? "Procesando..." : "Comprar ahora"}
                </Button>
              </div>
            ))}
          </div>

          <div className="bg-secondary/30 rounded-lg p-4 text-sm text-muted-foreground">
            <p className="flex items-center gap-2">
              <CheckCircle2 className="w-4 h-4 text-accent" />
              Los créditos se agregan automáticamente después del pago confirmado
            </p>
          </div>
        </div>

        {/* Payment History */}
        {payments.length > 0 && (
          <div className="bg-card border border-border rounded-xl p-6 space-y-6">
            <div>
              <h2 className="text-2xl font-bold text-foreground mb-2">Historial de pagos</h2>
              <p className="text-muted-foreground">Tus compras de créditos</p>
            </div>

            <Separator />

            <div className="space-y-3">
              {payments.slice(0, 10).map((payment: any) => (
                <div
                  key={payment.id}
                  className="flex items-center justify-between p-4 rounded-lg hover:bg-secondary/50 transition-colors"
                >
                  <div className="flex items-center gap-4 flex-1 min-w-0">
                    <div className={`w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0 ${
                      payment.status === 'approved' ? 'bg-green-500/10' :
                      payment.status === 'pending' ? 'bg-orange-500/10' : 'bg-red-500/10'
                    }`}>
                      <CheckCircle2 className={`w-5 h-5 ${
                        payment.status === 'approved' ? 'text-green-500' :
                        payment.status === 'pending' ? 'text-orange-500' : 'text-red-500'
                      }`} />
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-foreground">{payment.package_name}</p>
                      <p className="text-xs text-muted-foreground">
                        {format(new Date(payment.created_at), "dd MMM yyyy, HH:mm", { locale: es })}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-4">
                    <div className="text-right">
                      <p className="text-sm font-semibold text-foreground">
                        ${payment.amount.toLocaleString('es-AR')}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        +{payment.credits_amount} min
                      </p>
                    </div>
                    <div className={`px-2 py-1 rounded text-xs font-medium ${
                      payment.status === 'approved' ? 'bg-green-500/10 text-green-500' :
                      payment.status === 'pending' ? 'bg-orange-500/10 text-orange-500' :
                      'bg-red-500/10 text-red-500'
                    }`}>
                      {payment.status === 'approved' ? 'Aprobado' :
                       payment.status === 'pending' ? 'Pendiente' : 'Rechazado'}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Usage History */}
        <div className="bg-card border border-border rounded-xl p-6 space-y-6">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-2xl font-bold text-foreground mb-2">Historial de uso</h2>
              <p className="text-muted-foreground">Últimas transcripciones realizadas</p>
            </div>
          </div>

          <Separator />

          <div className="space-y-3">
            {completedTranscriptions.length > 0 ? (
              completedTranscriptions.slice(0, 10).map((item: any) => (
                <div
                  key={item.id}
                  className="flex items-center justify-between p-4 rounded-lg hover:bg-secondary/50 transition-colors"
                >
                  <div className="flex items-center gap-4 flex-1 min-w-0">
                    <div className="w-10 h-10 bg-accent/10 rounded-lg flex items-center justify-center flex-shrink-0">
                      <Clock className="w-5 h-5 text-accent" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-foreground truncate">{item.file_name}</p>
                      <p className="text-xs text-muted-foreground">
                        {format(new Date(item.created_at), "dd MMM yyyy, HH:mm", { locale: es })}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-4">
                    <span className="text-sm font-semibold text-foreground">-{item.credits_used} min</span>
                    <div className="w-2 h-2 bg-green-500 rounded-full" />
                  </div>
                </div>
              ))
            ) : (
              <div className="text-center py-8">
                <p className="text-muted-foreground">No hay transcripciones completadas aún</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
};

export default Credits;
