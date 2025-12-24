import * as React from 'react';

const PlayerDetail = ({ title, children }) => (
  <div className="mb-4 break-inside-avoid rounded-lg border border-white/10 bg-white/5 px-4 py-3">
    <h3 className="text-sm font-semibold uppercase tracking-wide text-gray-400 mb-3">
      {title}
    </h3>
    <div className="space-y-2">{children}</div>
  </div>
);

PlayerDetail.Item = function PlayerDetailItem({ label, children }) {
  return (
    <div>
      <p className="text-xs uppercase text-gray-400">{label}</p>
      <p className="text-white">{children}</p>
    </div>
  );
};

export default PlayerDetail;

