package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var (
	memRequestRatio, memLimitRatio, cpuRequestRatio, cpuLimitRatio float64
	allNamespaces                                                  bool
	namespace                                                      string
	logLevel                                                       string
)

type Options struct {
	memRequestRatio float64
	memLimitRatio   float64
	cpuRequestRatio float64
	cpuLimitRatio   float64
	allNamespaces   bool
	namespace       string
	logLevel        string
	streams         genericclioptions.IOStreams
	duration        int

	clientset        *kubernetes.Clientset
	metricsclientset *metricsv.Clientset
}

func newCmd(streams genericclioptions.IOStreams) *cobra.Command {
	opts := &Options{
		streams: streams,
	}
	cmd := &cobra.Command{
		Use:   "resource-report",
		Short: "resource-report",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(opts)
		},
	}
	cmd.PersistentFlags().BoolVarP(&opts.allNamespaces, "all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces.")
	cmd.PersistentFlags().StringVarP(&opts.namespace, "namespace", "n", "", "Container to search")
	cmd.MarkFlagsMutuallyExclusive("all-namespaces", "namespace")

	cmd.PersistentFlags().StringVarP(&opts.logLevel, "log-level", "l", "", "log level (trace|debug|info|warn|error|fatal|panic)")

	cmd.PersistentFlags().IntVarP(&opts.duration, "duration", "d", 1, "time to collect in seconds")

	cmd.PersistentFlags().Float64Var(&opts.memRequestRatio, "memory-request", 2.5, "memory request/usage ratio")
	cmd.PersistentFlags().Float64Var(&opts.memLimitRatio, "memory-limit", 1.5, "memory limit/request ration")
	cmd.PersistentFlags().Float64Var(&opts.cpuRequestRatio, "cpu-request", 100.0, "cpu request/usage ratio")
	cmd.PersistentFlags().Float64Var(&opts.cpuLimitRatio, "cpu-limit", 10.0, "cpu limit/request ratio")
	return cmd
}

func run(opts *Options) error {
	// progress
	pw := progress.NewWriter()
	pw.SetAutoStop(true)
	pw.SetMessageWidth(24)
	pw.Style().Visibility.Time = false
	go pw.Render()

	setup(opts, pw)

	specs := PodContainers{}
	pods(opts, pw, specs)
	metrics(opts, pw, specs)

	for pw.IsRenderInProgress() {
		time.Sleep(100 * time.Millisecond)
	}

	render(specs)

	return nil
}

func setup(opts *Options, pw progress.Writer) error {
	setupTracker := progress.Tracker{Message: "Setup"}
	pw.AppendTracker(&setupTracker)

	configFlags := genericclioptions.NewConfigFlags(true)
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		panic(err)
	}

	clientcfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		nil,
	)

	if opts.namespace == "" && !opts.allNamespaces {
		namespace, _, err := clientcfg.Namespace()
		opts.namespace = namespace
		if err != nil {
			panic(err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	opts.clientset = clientset

	metricsclientset, err := metricsv.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	opts.metricsclientset = metricsclientset

	setupTracker.MarkAsDone()
	return nil
}

func pods(opts *Options, pw progress.Writer, specs PodContainers) error {
	podTracker := progress.Tracker{Message: "Collecting pods", Units: progress.UnitsDefault}
	pw.AppendTracker(&podTracker)

	pods, err := opts.clientset.CoreV1().Pods(opts.namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, pod := range pods.Items {
		if pod.Status.Phase != "Running" {
			continue
		}
		for _, container := range pod.Spec.Containers {
			id := fmt.Sprintf("%s__%s__%s", pod.Namespace, pod.Name, container.Name)
			specs[id] = PodContainer{
				Namespace:     pod.Namespace,
				PodName:       pod.Name,
				ContainerName: container.Name,
				Requests: Resources{
					CPU:    uint64(container.Resources.Requests.Cpu().MilliValue()),
					Memory: container.Resources.Requests.Memory().ToDec().AsDec().UnscaledBig().Uint64(),
				},
				Limits: Resources{
					CPU:    uint64(container.Resources.Limits.Cpu().MilliValue()),
					Memory: container.Resources.Limits.Memory().ToDec().AsDec().UnscaledBig().Uint64(),
				},
			}
		}
	}
	podTracker.MarkAsDone()
	return nil
}

func metrics(opts *Options, pw progress.Writer, specs PodContainers) error {

	loops := int(math.Ceil(float64(opts.duration) / 15))
	metricTracker := progress.Tracker{Message: "Collecting metrics", Total: int64(loops), Units: progress.UnitsDefault}
	pw.AppendTracker(&metricTracker)
	for i := 1; i <= loops; i += 1 {
		metrics, err := opts.metricsclientset.MetricsV1beta1().PodMetricses(opts.namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		for _, metric := range metrics.Items {
			for _, container := range metric.Containers {
				id := fmt.Sprintf("%s__%s__%s", metric.Namespace, metric.Name, container.Name)
				if s, ok := specs[id]; ok {
					s.Usage = Resources{
						CPU:    max(s.Usage.CPU, uint64(container.Usage.Cpu().MilliValue())),
						Memory: max(s.Usage.Memory, container.Usage.Memory().ToDec().AsDec().UnscaledBig().Uint64()),
					}
					specs[id] = s
				}
			}
		}
		metricTracker.Increment(1)
		if i != loops {
			time.Sleep(15 * time.Second)

		}
	}
	metricTracker.MarkAsDone()
	return nil
}

func render(specs PodContainers) {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.Style().Format.Footer = text.FormatDefault
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 3, Align: text.AlignRight},
		{Number: 4, Align: text.AlignRight},
		{Number: 5, Align: text.AlignRight},
		{Number: 6, Align: text.AlignRight},
		{Number: 7, Align: text.AlignRight},
		{Number: 8, Align: text.AlignRight},
	})

	rows := table.Row{"Pod", "Container", "Req CPU", "Lim CPU", "Use CPU", "Req Mem", "Lim Mem", "Use Mem"}
	t.AppendHeader(rows)

	t.SuppressEmptyColumns()
	cpursum := uint64(0)
	cpuusum := uint64(0)
	cpulsum := uint64(0)
	memrsum := uint64(0)
	memusum := uint64(0)
	memlsum := uint64(0)
	for _, s := range specs.Sorted() {
		data := []interface{}{
			s.PodName,
			s.ContainerName,
		}
		cpursum += s.Requests.CPU
		cpulsum += s.Limits.CPU
		cpuusum += s.Usage.CPU
		memrsum += s.Requests.Memory
		memlsum += s.Limits.Memory
		memusum += s.Usage.Memory
		data = append(data, tuning(s.Requests.CPU, s.Limits.CPU, s.Usage.CPU, s.Requests.Memory, s.Limits.Memory, s.Usage.Memory)...)
		t.AppendRow(data)
	}
	footerRow := table.Row{"Total", "", clean(fmt.Sprintf("%dm", cpursum)), clean(fmt.Sprintf("%dm", cpulsum)), clean(fmt.Sprintf("%dm", cpuusum)), clean(humanize.IBytes(memrsum)), clean(humanize.IBytes(memlsum)), clean(humanize.IBytes(memusum))}
	t.AppendFooter(footerRow)
	t.AppendFooter(table.Row{""})
	t.AppendFooter(table.Row{
		text.FgHiRed.Sprint("low") + " | " + text.FgHiCyan.Sprint("high") + " | " + text.FgBlue.Sprint("unset"),
	})
	t.Render()
}

func main() {
	root := newCmd(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func tuning(cpur, cpul, cpuu, memr, meml, memu uint64) []interface{} {
	scpur := clean(fmt.Sprintf("%dm", cpur))
	scpul := clean(fmt.Sprintf("%dm", cpul))
	scpuu := clean(fmt.Sprintf("%dm", cpuu))
	smemr := clean(humanize.IBytes(memr))
	smeml := clean(humanize.IBytes(meml))
	smemu := clean(humanize.IBytes(memu))

	// memory
	if memr > uint64(float64(memu)*memRequestRatio) {
		smemr = text.FgHiCyan.Sprint(smemr)
		// smemu = text.FgYellow.Sprint(smemu)
	}
	if meml > uint64(float64(memr)*memLimitRatio) {
		smeml = text.FgHiCyan.Sprint(smeml)
	}

	if memr < memu {
		smemr = text.FgHiRed.Sprint(smemr)
	}
	if float64(meml)*0.9 <= float64(memu) {
		smeml = text.FgHiMagenta.Sprint(smeml)
	}

	// cpu
	if cpur > uint64(float64(cpuu+1)*cpuRequestRatio) {
		scpur = text.FgHiCyan.Sprint(scpur)
		// scpuu = text.FgYellow.Sprint(scpuu)
	}
	if cpul > uint64(float64(cpur)*cpuLimitRatio) {
		scpul = text.FgHiCyan.Sprint(scpul)
	}

	if cpur < cpuu {
		scpur = text.FgHiRed.Sprint(scpur)
	}
	if float64(cpul)*0.9 <= float64(cpuu) {
		scpul = text.FgHiMagenta.Sprint(scpul)
	}

	return []interface{}{
		highlightNotSet(cpur, scpur),
		highlightNotSet(cpul, scpul),
		scpuu,
		highlightNotSet(memr, smemr),
		highlightNotSet(meml, smeml),
		smemu,
	}
}

func highlightNotSet(i uint64, s string) string {
	if i == 0 {
		return text.FgBlue.Sprint("null")
	}
	return s
}

func clean(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

func truncateText(s string, max int) string {
	if max > len(s) {
		return s
	}
	return s[:strings.LastIndexAny(s[:max], " .,:;-")] + "..."
}

type PodContainers map[string]PodContainer

func (p PodContainers) Sorted() []PodContainer {
	keys := make([]string, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	v := make([]PodContainer, 0, len(p))
	for _, k := range keys {
		v = append(v, p[k])
	}

	return v
}

type PodContainer struct {
	Namespace     string
	PodName       string
	ContainerName string
	Requests      Resources
	Limits        Resources
	Usage         Resources
}
type Resources struct {
	CPU    uint64
	Memory uint64
}
