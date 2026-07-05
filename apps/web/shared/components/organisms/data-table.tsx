"use client";

import { Spinner } from "@/shared/components/atoms/spinner";

type DataTableProps = {
  headers: string[];
  rows: React.ReactNode[][];
  loading?: boolean;
  loadingLabel?: string;
  emptyTitle?: string;
  emptyDescription?: string;
};

export function DataTable({
  headers,
  rows,
  loading = false,
  loadingLabel = "Loading records...",
  emptyTitle = "No records found",
  emptyDescription = "Adjust search/filter or create a new record.",
}: DataTableProps) {
  if (loading) {
    return (
      <div className="flex items-center gap-2 p-4 text-sm text-muted-foreground">
        <Spinner />
        {loadingLabel}
      </div>
    );
  }

  if (rows.length === 0) {
    return (
      <div className="p-10 text-center">
        <p className="text-sm font-medium">{emptyTitle}</p>
        <p className="mt-1 text-sm text-muted-foreground">{emptyDescription}</p>
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-left text-sm">
        <thead className="bg-slate-50 text-xs uppercase text-muted-foreground">
          <tr>
            {headers.map((header) => (
              <th key={header} className="px-4 py-3">
                {header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, index) => (
            <tr key={index} className="border-t hover:bg-slate-50">
              {row.map((cell, cellIndex) => (
                <td key={cellIndex} className="px-4 py-3">
                  {cell}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
