const defaultDateFormatter = new Intl.DateTimeFormat('en-US', {
  month: 'short',
  day: 'numeric',
  year: '2-digit',
});

const numberFormatter = new Intl.NumberFormat('en-US');

export const formatDate = (value) => {
  if (!value) return '';
  const date = value instanceof Date ? value : new Date(value);
  if (Number.isNaN(date.getTime())) return '';
  return defaultDateFormatter.format(date);
};

export const formatNumber = (value) => {
  if (value == null) return '0';
  return numberFormatter.format(value);
};

export const formatCompactNumber = (value) => {
  if (value == null) return '0';
  const absValue = Math.abs(value);
  if (absValue >= 1_000_000_000) {
    return `${(value / 1_000_000_000).toFixed(1).replace(/\.0$/, '')}b`;
  }
  if (absValue >= 1_000_000) {
    return `${(value / 1_000_000).toFixed(1).replace(/\.0$/, '')}m`;
  }
  if (absValue >= 1_000) {
    return `${(value / 1_000).toFixed(1).replace(/\.0$/, '')}k`;
  }
  return numberFormatter.format(value);
};

