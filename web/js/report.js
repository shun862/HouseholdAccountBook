function initReportChart(labels, incomeData, expenseData) {
    // 最大値
    const maxValue = Math.max(...incomeData, ...expenseData);
    const ceilMaxValue = CeilMaxValue(maxValue);
    // 中間値
    const stepValue = ceilMaxValue / 2;

    // 選択月インデックス
    let currentMonthIndex = labels.length - 1;

    // 選択月強調表示
    const highlightPlugin = {
        id: "highlightCurrentMonth",
        beforeDraw(chart) {
            if (currentMonthIndex === null) {
                return;
            }

            const {
                ctx,
                chartArea,
                scales
            } = chart;

            const xScale = scales.x;
            const center = xScale.getPixelForValue(
                currentMonthIndex
            );
            const width = xScale.width / labels.length;

            ctx.save();

            ctx.fillStyle = "rgba(37, 99, 235, 0.08)";
            ctx.fillRect(
                center - width / 2,
                chartArea.top - 20,
                width,
                chartArea.bottom - chartArea.top + 50
            );

            ctx.restore();

            // 集計表示
            const totalTitle = document.getElementById("total-htitle");
            const incomeAmount = document.getElementById("income-amount")
            const expenseAmount = document.getElementById("expense-amount")
            const balanceAmount = document.getElementById("balance-amount")

            const label = labels[currentMonthIndex];
            const income = incomeData[currentMonthIndex];
            const expense = expenseData[currentMonthIndex];
            totalTitle.textContent = label + "の集計";
            incomeAmount.innerHTML = `${income.toLocaleString()}<span>円</span>`;
            expenseAmount.innerHTML = `${expense.toLocaleString()}<span>円</span>`;
            balanceAmount.innerHTML = `${(income - expense).toLocaleString()}<span>円</span>`;
        }
    };
    // 凡例の余白
    const legendMarginPlugin = {
        id: 'legendMargin',
        beforeInit(chart) {
            const originalFit = chart.legend.fit;
            chart.legend.fit = function () {
                originalFit.bind(chart.legend)();
                this.height += 30;
            };
        }
    };

    const chart = new Chart(document.getElementById("report-chart"), {
        type: "bar",
        data: {
            labels: labels,
            datasets: [
                {
                    label: "収入",
                    data: incomeData,
                    backgroundColor: "#2563eb",
                    borderRadius: 5,
                    categoryPercentage: 0.7,
                    barPercentage: 0.8
                },
                {
                    label: "支出",
                    data: expenseData,
                    backgroundColor: "#ef4444",
                    categoryPercentage: 0.7,
                    barPercentage: 0.8
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        stepSize: stepValue
                    }
                },
                x: {
                    grid: {
                        display: false
                    }
                }
            },
            plugins: {
                legend: {
                    position: "top",
                    labels: {
                        boxWidth: 20,
                        boxHeight: 12,
                        padding: 30,
                        font: {
                            size: 12
                        }
                    }
                }
            },
        },
        plugins: [
            highlightPlugin,
            legendMarginPlugin
        ],
    });

    // 月選択イベント
    document.getElementById("report-chart").onclick = function (event) {
        const points = chart.getElementsAtEventForMode(
            event,
            "index",
            {
                intersect: false
            },
            true
        );

        if (points.length > 0) {
            currentMonthIndex = points[0].index;
            chart.update();
        }
    };
}

function CeilMaxValue(maxValue) {
    if (maxValue <= 10000) return 10000;
    return Math.ceil(maxValue / 10000) * 10000;
}