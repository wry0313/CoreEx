import { SolidApexCharts } from "solid-apexcharts";
import { createSignal } from "solid-js";
import { list } from "./options";

export default function CandleGraph() {
  const [options] = createSignal({
    chart: {
      id: "solidchart-example",
    },
  });
  const [series] = createSignal(list);

  return (
    <SolidApexCharts
      width="500"
      type="candlestick"
      options={options()}
      series={series()}
    />
  );
}