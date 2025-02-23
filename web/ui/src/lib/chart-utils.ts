import type { MalakIntegrationDataPoint } from "@/client/Api";

export interface ChartDataPoint {
  name: string;
  value: number;
}

export const formatChartData = (dataPoints: MalakIntegrationDataPoint[] | undefined): ChartDataPoint[] => {
  if (!dataPoints) return [];

  return dataPoints.map(point => {
    const value = point.data_point_type === 'currency'
      ? (point.point_value || 0) / 100
      : point.point_value || 0;

    return {
      name: point.point_name || '',
      value,
    };
  });
};

export const formatTooltipValue = (value: number, dataPointType?: string): [string | number, string] => {
  if (dataPointType === 'currency') {
    return [`$${value.toFixed(2)}`, 'Value'];
  }
  return [value, 'Value'];
};

export const getChartColors = (index: number): string => {
  return `hsl(${index * 45}, 70%, 50%)`;
}; 