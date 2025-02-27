import {
  OptionsWithLegend,
  OptionsWithTextFormatting,
  OptionsWithTooltip,
  AxisConfig,
  BarValueVisibility,
  GraphGradientMode,
  HideableFieldConfig,
  StackingMode,
} from '@grafana/schema';
import { VizOrientation } from '@grafana/data';

/**
 * @alpha
 */
export interface BarChartOptions extends OptionsWithLegend, OptionsWithTooltip, OptionsWithTextFormatting {
  orientation: VizOrientation;
  stacking: StackingMode;
  showValue: BarValueVisibility;
  barWidth: number;
  groupWidth: number;
}

/**
 * @alpha
 */
export interface BarChartFieldConfig extends AxisConfig, HideableFieldConfig {
  lineWidth?: number; // 0
  fillOpacity?: number; // 100
  gradientMode?: GraphGradientMode;
}

/**
 * @alpha
 */
export const defaultBarChartFieldConfig: BarChartFieldConfig = {
  lineWidth: 1,
  fillOpacity: 80,
  gradientMode: GraphGradientMode.None,
  axisSoftMin: 0,
};
