import { MalakIntegrationDataPointType } from "@/client/Api";
import type { MalakIntegrationDataPoint } from "@/client/Api";

export interface ChartDataPoint {
  name: string;
  value: number;
}

export const formatChartData = (
  dataPoints: MalakIntegrationDataPoint[] | undefined,
  dataPointType?: MalakIntegrationDataPointType
): ChartDataPoint[] => {
  if (!dataPoints) return [];

  return dataPoints.map(point => {
    const value = dataPointType === MalakIntegrationDataPointType.IntegrationDataPointTypeCurrency
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
    return [`$${value.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`, 'Value'];
  }
  return [value, 'Value'];
};

export const getChartColors = (index: number): string => {
  return `hsl(${index * 45}, 70%, 50%)`;
}; 