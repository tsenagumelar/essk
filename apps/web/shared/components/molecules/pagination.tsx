"use client";

import { Button } from "@/shared/components/atoms/button";
import { SelectField } from "@/shared/components/atoms/select-field";

const defaultPageSizeOptions = [5, 10, 20];

type PaginationProps = {
  page: number;
  pageSize: number;
  totalPages: number;
  totalItems: number;
  firstRecord: number;
  lastRecord: number;
  pageSizeOptions?: number[];
  onPageChange: (page: number) => void;
  onPageSizeChange: (pageSize: number) => void;
};

export function Pagination({
  page,
  pageSize,
  totalPages,
  totalItems,
  firstRecord,
  lastRecord,
  pageSizeOptions = defaultPageSizeOptions,
  onPageChange,
  onPageSizeChange,
}: PaginationProps) {
  return (
    <div className="flex flex-col gap-3 border-t px-4 py-3 text-sm text-muted-foreground md:flex-row md:items-center md:justify-between">
      <p>
        Showing {firstRecord}-{lastRecord} of {totalItems}
      </p>
      <div className="flex items-center gap-2">
        <SelectField
          className="h-9 w-auto"
          value={String(pageSize)}
          options={pageSizeOptions.map((size) => ({ value: String(size), label: `${size} / page` }))}
          onChange={(event) => onPageSizeChange(Number(event.target.value))}
        />
        <Button type="button" variant="outline" size="sm" disabled={page <= 1} onClick={() => onPageChange(Math.max(1, page - 1))}>
          Prev
        </Button>
        <span className="px-2">
          Page {page} of {totalPages}
        </span>
        <Button type="button" variant="outline" size="sm" disabled={page >= totalPages} onClick={() => onPageChange(Math.min(totalPages, page + 1))}>
          Next
        </Button>
      </div>
    </div>
  );
}
