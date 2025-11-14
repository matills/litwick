import { FileText, CheckCircle2, Clock, Flame } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api";

export const StatsCards = () => {
  const { data } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: apiClient.getDashboard,
  });

  const stats = [
    {
      icon: FileText,
      label: "Total",
      value: data?.stats?.total_transcriptions || 0,
      color: "text-accent"
    },
    {
      icon: CheckCircle2,
      label: "Completadas",
      value: data?.stats?.completed_count || 0,
      color: "text-green-500"
    },
    {
      icon: Clock,
      label: "Procesando",
      value: data?.stats?.processing_count || 0,
      color: "text-orange-500"
    },
    {
      icon: Flame,
      label: "Créditos",
      value: `${data?.stats?.credits_remaining || 0} min`,
      color: "text-orange-600"
    }
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
      {stats.map((stat, index) => {
        const Icon = stat.icon;
        return (
          <div
            key={index}
            className="p-6 bg-card border border-border rounded-xl hover:shadow-soft transition-all"
          >
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-secondary rounded-lg flex items-center justify-center">
                <Icon className={`w-6 h-6 ${stat.color}`} />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">{stat.label}</p>
                <p className="text-3xl font-bold text-foreground">{stat.value}</p>
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};
