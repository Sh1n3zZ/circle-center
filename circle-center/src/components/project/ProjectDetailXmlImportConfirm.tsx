import type {
  ConfirmImportRequest,
  IconImportComponent,
} from '@/api/manager/types';
import { xmlioApi } from '@/api/manager/xmlio';
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { Button } from '@/components/ui/button';
import { useState } from 'react';
import { toast } from 'sonner';
import ConfirmTable from './ProjectDetailXmlImportConfirmTable';

export default function ProjectDetailXmlImportConfirm({
  projectId,
  initial,
  onClose,
}: {
  projectId: number;
  initial: IconImportComponent[];
  onClose: () => void;
}) {
  const [data, setData] = useState<IconImportComponent[]>(initial);
  const [loading, setLoading] = useState(false);
  const [isDialogOpen, setIsDialogOpen] = useState(true);

  const handleConfirm = async () => {
    try {
      setLoading(true);
      const body: ConfirmImportRequest = { projectId, components: data };
      const res = await xmlioApi.confirmImport(body);
      if (res.status !== 'success')
        throw new Error(res.message || 'Import failed');
      toast.success(
        `Import completed: ${res.summary?.created} created, ${res.summary?.duplicates} duplicates, ${res.summary?.errors} errors`
      );
      setIsDialogOpen(false);
      onClose();
    } catch (e: any) {
      toast.error(e?.response?.data?.message || e.message || 'Import failed');
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setIsDialogOpen(false);
    onClose();
  };

  return (
    <AlertDialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
      <AlertDialogContent className='!max-w-[95vw] !w-[95vw] max-h-[90vh] flex flex-col sm:!max-w-[95vw] lg:!max-w-[90vw] xl:!max-w-[85vw]'>
        <AlertDialogHeader>
          <AlertDialogTitle>Confirm XML Import</AlertDialogTitle>
          <AlertDialogDescription>
            Review and edit the parsed components before importing. Found{' '}
            {data.length} components to import.
          </AlertDialogDescription>
        </AlertDialogHeader>

        {/* Scrollable content area */}
        <div className='flex-1 overflow-auto min-h-0'>
          <ConfirmTable data={data} onChange={setData} />
        </div>

        <AlertDialogFooter className='flex-shrink-0'>
          <Button variant='outline' onClick={handleClose} disabled={loading}>
            Cancel
          </Button>
          <Button onClick={handleConfirm} disabled={loading}>
            {loading ? 'Importing...' : 'Confirm Import'}
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
