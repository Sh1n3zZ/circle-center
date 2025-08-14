import { useEffect, useState } from "react";
import { tokenApi } from "@/api/manager/token";
import type { ProjectTokenModel } from "@/api/manager/types";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";

export default function ProjectDetailTokenList({ projectId, refreshSignal = 0 }: { projectId: number; refreshSignal?: number }) {
  const [tokens, setTokens] = useState<ProjectTokenModel[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchTokens = async () => {
    try {
      setLoading(true);
      const res = await tokenApi.list(projectId);
      setTokens(res.data || []);
    } catch (e: any) {
      toast.error(e?.response?.data?.message || e.message || "Failed to load tokens");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTokens();
  }, [projectId, refreshSignal]);

  const handleDelete = async (id: number) => {
    try {
      await tokenApi.delete(projectId, id);
      toast.success("Token deleted");
      setTokens((prev) => prev.filter((t) => t.id !== id));
    } catch (e: any) {
      toast.error(e?.response?.data?.message || e.message || "Delete failed");
    }
  };

  return (
    <div className="space-y-3">
      {loading ? (
        <div className="text-sm text-muted-foreground">Loading...</div>
      ) : tokens.length === 0 ? (
        <div className="text-sm text-muted-foreground">No tokens yet</div>
      ) : (
        <div className="space-y-2">
          {tokens.map((t) => (
            <div key={t.id} className="border rounded-md p-3 flex items-center justify-between">
              <div className="space-y-0.5">
                <div className="font-medium">{t.name}</div>
                <div className="text-xs text-muted-foreground">
                  Active: {t.active ? "Yes" : "No"} · Created: {new Date(t.created_at).toLocaleString()} · Last used: {t.last_used_at ? new Date(t.last_used_at).toLocaleString() : "Never"}
                </div>
              </div>
              <div className="flex gap-2">
                <Button variant="destructive" size="sm" onClick={() => handleDelete(t.id)}>Delete</Button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}


