import { Toggle } from '@base-ui/react/toggle';
import { ToggleGroup } from '@base-ui/react/toggle-group';

const DEFAULT_TIME_RANGES = ['1w', '1m', '1y', 'all'];

export default function TimeRangeToggle({
  value = '1w',
  onChange = () => {},
  options = DEFAULT_TIME_RANGES,
  className,
}) {
  return (
    <ToggleGroup
      type="single"
      value={value}
      onValueChange={onChange}
      aria-label="Select a time range"
      className={[
        'inline-flex items-center gap-1 rounded-lg border border-white/20 bg-white/5 p-1 text-xs font-semibold',
        className,
      ]
        .filter(Boolean)
        .join(' ')}
    >
      {options.map((range) => (
        <Toggle
          key={range}
          value={range}
          className="
            rounded-sm px-3 py-1 text-xs font-semibold text-white transition
            aria-[pressed=true]:bg-zinc-800
            aria-[pressed=false]:bg-transparent
            aria-[pressed=false]:hover:bg-zinc-800
          "
        >
          {range}
        </Toggle>
      ))}
    </ToggleGroup>
  );
}
