import * as React from 'react';
import { Select as BaseSelect } from '@base-ui/react/select';

export default function Select({
  items,
  value,
  onValueChange,
  placeholder = 'Select...',
  className = '',
  triggerClassName = '',
  ...props
}) {
  const selectedItem = items.find((item) => item.value === value);

  return (
    <BaseSelect.Root value={value} onValueChange={onValueChange} {...props}>
      <BaseSelect.Trigger
        className={`bg-zinc-900 border border-white/15 rounded-lg px-3 py-1.5 text-sm focus:border-blue-400 focus:outline-none flex items-center min-w-28 ${triggerClassName}`}
      >
        <div className="flex-1 text-left">
          {selectedItem ? selectedItem.label : placeholder}
        </div>
        <BaseSelect.Icon className="text-gray-400">
          <ChevronDownIcon />
        </BaseSelect.Icon>
      </BaseSelect.Trigger>
      <BaseSelect.Portal>
        <BaseSelect.Positioner className="absolute z-50 mt-1" sideOffset={4}>
          <BaseSelect.Popup
            className={`bg-zinc-900 border border-white/15 rounded-lg shadow-lg w-full min-w-28 ${className}`}
          >
            <BaseSelect.List className="py-1">
              {items.map(({ label, value: itemValue }) => (
                <BaseSelect.Item
                  key={itemValue}
                  value={itemValue}
                  className="px-3 py-2 hover:bg-white/10 cursor-pointer text-white text-sm flex items-center"
                >
                  <BaseSelect.ItemIndicator className="mr-2 flex items-center justify-center">
                    <CheckIcon />
                  </BaseSelect.ItemIndicator>
                  <BaseSelect.ItemText>{label}</BaseSelect.ItemText>
                </BaseSelect.Item>
              ))}
            </BaseSelect.List>
          </BaseSelect.Popup>
        </BaseSelect.Positioner>
      </BaseSelect.Portal>
    </BaseSelect.Root>
  );
}

function ChevronDownIcon(props) {
  return (
    <svg
      width="12"
      height="12"
      viewBox="0 0 12 12"
      fill="none"
      stroke="currentcolor"
      strokeWidth="2"
      {...props}
    >
      <path d="M3 4.5L6 7.5L9 4.5" />
    </svg>
  );
}

function CheckIcon(props) {
  return (
    <svg
      fill="currentcolor"
      width="10"
      height="10"
      viewBox="0 0 10 10"
      {...props}
    >
      <path d="M9.1603 1.12218C9.50684 1.34873 9.60427 1.81354 9.37792 2.16038L5.13603 8.66012C5.01614 8.8438 4.82192 8.96576 4.60451 8.99384C4.3871 9.02194 4.1683 8.95335 4.00574 8.80615L1.24664 6.30769C0.939709 6.02975 0.916013 5.55541 1.19372 5.24822C1.47142 4.94102 1.94536 4.91731 2.2523 5.19524L4.36085 7.10461L8.12299 1.33999C8.34934 0.993152 8.81376 0.895638 9.1603 1.12218Z" />
    </svg>
  );
}
