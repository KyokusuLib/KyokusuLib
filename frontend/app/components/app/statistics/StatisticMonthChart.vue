<script setup lang="ts">
import { computed } from 'vue'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from 'chart.js'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler,
)

interface Dataset {
  label: string
  color: string
}

interface DataPoint {
  label: string
  values: number[]
}

const props = defineProps<{
  datasets: Dataset[]
  data: DataPoint[]
  height?: number
}>()

const chartData = computed(() => ({
  labels: props.data.map((d) => d.label),
  datasets: props.datasets.map((ds, i) => ({
    label: ds.label,
    data: props.data.map((d) => d.values[i] ?? 0),
    borderColor: ds.color,
    backgroundColor: ds.color + '15',
    fill: true,
    tension: 0.3,
    pointRadius: 3,
    pointHoverRadius: 6,
    pointBackgroundColor: ds.color,
    pointBorderColor: '#fff',
    pointBorderWidth: 2,
    borderWidth: 2,
  })),
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: true,
      position: 'top' as const,
      align: 'top' as const,
      labels: {
        usePointStyle: true,
        boxWidth: 8,
        boxHeight: 8,
        padding: 16,
        color: '#a1a1aa',
        font: { size: 12 },
      },
    },
    tooltip: {
      mode: 'index' as const,
      intersect: false,
      backgroundColor: 'rgba(24, 24, 27, 0.95)',
      titleColor: '#e4e4e7',
      bodyColor: '#a1a1aa',
      borderColor: 'rgba(63, 63, 70, 0.5)',
      borderWidth: 1,
      padding: 12,
      cornerRadius: 12,
    },
  },
  scales: {
    x: {
      grid: { display: false },
      ticks: {
        color: '#a1a1aa',
        font: { size: 11 },
      },
    },
    y: {
      beginAtZero: true,
      grid: {
        color: 'rgba(113, 113, 122, 0.15)',
        drawTicks: false,
      },
      border: { display: false },
      ticks: {
        color: '#a1a1aa',
        font: { size: 11 },
        padding: 8,
      },
    },
  },
  interaction: {
    mode: 'index' as const,
    intersect: false,
  },
}
</script>

<template>
  <div class="flex w-full gap-4 py-4 h-120">
    <Line :data="chartData" :options="chartOptions" />
  </div>
</template>
