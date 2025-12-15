import { DashboardLayout } from "@/components/DashboardLayout";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { User, Bell, Globe, AlertTriangle } from "lucide-react";
import { Separator } from "@/components/ui/separator";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api";
import { supabase } from "@/lib/supabase";
import { useState, useEffect } from "react";
import { toast } from "sonner";

const Settings = () => {
  const [email, setEmail] = useState("");
  const queryClient = useQueryClient();

  const [defaultLanguage, setDefaultLanguage] = useState("es");
  const [defaultExportFormat, setDefaultExportFormat] = useState("srt");
  const [includeTimestamps, setIncludeTimestamps] = useState(true);
  const [detectSpeakers, setDetectSpeakers] = useState(true);
  const [emailNotifications, setEmailNotifications] = useState(true);
  const [promotionalEmails, setPromotionalEmails] = useState(false);

  const { data: dashboardData } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: apiClient.getDashboard,
  });

  useEffect(() => {
    supabase.auth.getSession().then(({ data: { session } }) => {
      if (session?.user?.email) {
        setEmail(session.user.email);
      }
    });

    if (dashboardData?.user) {
      setDefaultLanguage(dashboardData.user.default_language || "es");
      setDefaultExportFormat(dashboardData.user.default_export_format || "srt");
      setIncludeTimestamps(dashboardData.user.include_timestamps ?? true);
      setDetectSpeakers(dashboardData.user.detect_speakers ?? true);
      setEmailNotifications(dashboardData.user.email_notifications ?? true);
      setPromotionalEmails(dashboardData.user.promotional_emails ?? false);
    }
  }, [dashboardData]);

  const updateSettingsMutation = useMutation({
    mutationFn: apiClient.updateSettings,
    onSuccess: () => {
      toast.success("Configuración guardada correctamente");
      queryClient.invalidateQueries({ queryKey: ['dashboard-stats'] });
    },
    onError: () => {
      toast.error("Error al guardar la configuración");
    },
  });

  const handleSaveSettings = async () => {
    updateSettingsMutation.mutate({
      default_language: defaultLanguage,
      default_export_format: defaultExportFormat,
      include_timestamps: includeTimestamps,
      detect_speakers: detectSpeakers,
      email_notifications: emailNotifications,
      promotional_emails: promotionalEmails,
    });
  };

  return (
    <DashboardLayout>
      <div className="max-w-4xl mx-auto space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-foreground mb-2">Configuración</h1>
          <p className="text-muted-foreground">
            Administra tu cuenta y preferencias de la aplicación
          </p>
        </div>

        <div className="bg-card border border-border rounded-xl p-6 space-y-6">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-accent/10 rounded-lg flex items-center justify-center">
              <User className="w-5 h-5 text-accent" />
            </div>
            <div>
              <h2 className="text-xl font-semibold text-foreground">Cuenta</h2>
              <p className="text-sm text-muted-foreground">Información personal y credenciales</p>
            </div>
          </div>

          <Separator />

          <div className="space-y-4">
            <div className="grid gap-2">
              <Label htmlFor="email">Correo electrónico</Label>
              <Input
                id="email"
                type="email"
                placeholder="tu@email.com"
                value={email}
                disabled
                className="bg-muted"
              />
              <p className="text-xs text-muted-foreground">
                El email está vinculado a tu cuenta de Supabase
              </p>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="plan">Plan actual</Label>
              <Input
                id="plan"
                value={dashboardData?.user?.plan || "free"}
                disabled
                className="bg-muted capitalize"
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="credits">Créditos restantes</Label>
              <Input
                id="credits"
                value={`${dashboardData?.stats?.credits_remaining || 0} minutos`}
                disabled
                className="bg-muted"
              />
            </div>

            <div className="pt-4 border-t">
              <p className="text-sm text-muted-foreground mb-4">
                Para cambiar tu contraseña o actualizar tu email, por favor usa las opciones de Supabase Auth.
              </p>
            </div>
          </div>
        </div>

        <div className="bg-card border border-border rounded-xl p-6 space-y-6">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-accent/10 rounded-lg flex items-center justify-center">
              <Globe className="w-5 h-5 text-accent" />
            </div>
            <div>
              <h2 className="text-xl font-semibold text-foreground">Preferencias de Transcripción</h2>
              <p className="text-sm text-muted-foreground">Configuración predeterminada para nuevas transcripciones</p>
            </div>
          </div>

          <Separator />

          <div className="space-y-4">
            <div className="grid gap-2">
              <Label htmlFor="language">Idioma predeterminado</Label>
              <Select value={defaultLanguage} onValueChange={setDefaultLanguage}>
                <SelectTrigger id="language">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="es">Español</SelectItem>
                  <SelectItem value="en">Inglés</SelectItem>
                  <SelectItem value="fr">Francés</SelectItem>
                  <SelectItem value="de">Alemán</SelectItem>
                  <SelectItem value="pt">Portugués</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="format">Formato de exportación predeterminado</Label>
              <Select value={defaultExportFormat} onValueChange={setDefaultExportFormat}>
                <SelectTrigger id="format">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="txt">Texto plano (.txt)</SelectItem>
                  <SelectItem value="srt">Subtítulos SRT (.srt)</SelectItem>
                  <SelectItem value="vtt">WebVTT (.vtt)</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="flex items-center justify-between py-2">
              <div className="space-y-0.5">
                <Label htmlFor="timestamps">Incluir timestamps</Label>
                <p className="text-sm text-muted-foreground">
                  Agregar marcas de tiempo en las exportaciones
                </p>
              </div>
              <Switch
                id="timestamps"
                checked={includeTimestamps}
                onCheckedChange={setIncludeTimestamps}
              />
            </div>

            <div className="flex items-center justify-between py-2">
              <div className="space-y-0.5">
                <Label htmlFor="speakers">Detectar múltiples speakers</Label>
                <p className="text-sm text-muted-foreground">
                  Identificar diferentes voces automáticamente
                </p>
              </div>
              <Switch
                id="speakers"
                checked={detectSpeakers}
                onCheckedChange={setDetectSpeakers}
              />
            </div>

            <Separator />

            <Button
              onClick={handleSaveSettings}
              disabled={updateSettingsMutation.isPending}
              className="w-full sm:w-auto"
            >
              {updateSettingsMutation.isPending ? "Guardando..." : "Guardar preferencias"}
            </Button>
          </div>
        </div>

        <div className="bg-card border border-border rounded-xl p-6 space-y-6">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-accent/10 rounded-lg flex items-center justify-center">
              <Bell className="w-5 h-5 text-accent" />
            </div>
            <div>
              <h2 className="text-xl font-semibold text-foreground">Notificaciones</h2>
              <p className="text-sm text-muted-foreground">Gestiona cómo quieres recibir actualizaciones</p>
            </div>
          </div>

          <Separator />

          <div className="space-y-4">
            <div className="flex items-center justify-between py-2">
              <div className="space-y-0.5">
                <Label htmlFor="email-notif">Notificaciones por email</Label>
                <p className="text-sm text-muted-foreground">
                  Recibir actualizaciones cuando se completen transcripciones
                </p>
              </div>
              <Switch
                id="email-notif"
                checked={emailNotifications}
                onCheckedChange={setEmailNotifications}
              />
            </div>

            <div className="flex items-center justify-between py-2">
              <div className="space-y-0.5">
                <Label htmlFor="promo-notif">Emails promocionales</Label>
                <p className="text-sm text-muted-foreground">
                  Recibir novedades y ofertas especiales
                </p>
              </div>
              <Switch
                id="promo-notif"
                checked={promotionalEmails}
                onCheckedChange={setPromotionalEmails}
              />
            </div>

            <Separator />

            <Button
              onClick={handleSaveSettings}
              disabled={updateSettingsMutation.isPending}
              className="w-full sm:w-auto"
            >
              {updateSettingsMutation.isPending ? "Guardando..." : "Guardar notificaciones"}
            </Button>
          </div>
        </div>

        <div className="bg-destructive/5 border border-destructive/20 rounded-xl p-6 space-y-6">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-destructive/10 rounded-lg flex items-center justify-center">
              <AlertTriangle className="w-5 h-5 text-destructive" />
            </div>
            <div>
              <h2 className="text-xl font-semibold text-foreground">Zona de peligro</h2>
              <p className="text-sm text-muted-foreground">Acciones irreversibles</p>
            </div>
          </div>

          <Separator />

          <div className="space-y-4">
            <div>
              <h3 className="font-medium text-foreground mb-2">Eliminar cuenta</h3>
              <p className="text-sm text-muted-foreground mb-4">
                Una vez eliminada tu cuenta, no hay vuelta atrás. Todos tus datos y transcripciones se eliminarán permanentemente.
              </p>
              <Button
                variant="destructive"
                className="w-full sm:w-auto"
                onClick={() => toast.error("Funcionalidad en desarrollo")}
              >
                Eliminar mi cuenta
              </Button>
            </div>
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
};

export default Settings;
