import type { IconImportComponent } from '@/api/manager/types';
import { xmlioApi } from '@/api/manager/xmlio';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Separator } from '@/components/ui/separator';
import { useCallback, useState } from 'react';
import { toast } from 'sonner';

export default function ProjectDetailXmlImportPanel({
  onParsed,
}: {
  onParsed: (components: IconImportComponent[]) => void;
}) {
  const [appfilter, setAppfilter] = useState<File | null>(null);
  const [appmap, setAppmap] = useState<File | null>(null);
  const [theme, setTheme] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);

  const handleDrop = useCallback(
    (e: React.DragEvent<HTMLDivElement>) => {
      e.preventDefault();
      const files = Array.from(e.dataTransfer.files || []);
      files.forEach(f => {
        const name = f.name.toLowerCase();
        if (name.includes('appfilter')) setAppfilter(f);
        else if (name.includes('appmap')) setAppmap(f);
        else if (name.includes('theme')) setTheme(f);
        else if (name.endsWith('.xml')) {
          // Fallback: try to guess by size or first matched slot
          if (!appfilter) setAppfilter(f);
          else if (!appmap) setAppmap(f);
          else if (!theme) setTheme(f);
        }
      });
    },
    [appfilter, appmap, theme]
  );

  const parseNow = async () => {
    if (!appfilter && !appmap && !theme) {
      toast.error('Please select or drag at least one XML file');
      return;
    }
    try {
      setLoading(true);
      const res = await xmlioApi.parseForm(
        appfilter ?? undefined,
        appmap ?? undefined,
        theme ?? undefined
      );
      if (res.status !== 'success' || !res.components) {
        throw new Error(res.message || 'Parse failed');
      }
      onParsed(res.components);
      toast.success(
        `Parse successful, found ${res.components.length} components`
      );
    } catch (e: any) {
      toast.error(e?.response?.data?.message || e.message || 'Parse failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className='space-y-3'>
      <div
        className='border-dashed border rounded-md p-4 text-center text-xs text-muted-foreground'
        onDragOver={e => e.preventDefault()}
        onDrop={handleDrop}
      >
        Drag appfilter.xml / appmap.xml / theme_resources.xml here, or use file
        inputs below
      </div>

      <div className='space-y-2'>
        <div className='space-y-1'>
          <div className='text-xs font-medium'>appfilter.xml</div>
          <Input
            type='file'
            accept='.xml'
            className='text-xs'
            onChange={e => setAppfilter(e.target.files?.[0] || null)}
          />
          {appfilter && (
            <div className='text-xs text-muted-foreground'>
              {appfilter.name}
            </div>
          )}
        </div>
        <div className='space-y-1'>
          <div className='text-xs font-medium'>appmap.xml</div>
          <Input
            type='file'
            accept='.xml'
            className='text-xs'
            onChange={e => setAppmap(e.target.files?.[0] || null)}
          />
          {appmap && (
            <div className='text-xs text-muted-foreground'>{appmap.name}</div>
          )}
        </div>
        <div className='space-y-1'>
          <div className='text-xs font-medium'>theme_resources.xml</div>
          <Input
            type='file'
            accept='.xml'
            className='text-xs'
            onChange={e => setTheme(e.target.files?.[0] || null)}
          />
          {theme && (
            <div className='text-xs text-muted-foreground'>{theme.name}</div>
          )}
        </div>
      </div>

      <Separator />
      <div className='flex justify-end'>
        <Button
          onClick={parseNow}
          disabled={loading}
          size='sm'
          className='text-xs'
        >
          {loading ? 'Parsing...' : 'Parse & Preview'}
        </Button>
      </div>
    </div>
  );
}
