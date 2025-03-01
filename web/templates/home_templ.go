// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Home() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div class=\"max-w-2xl mx-auto\"><div class=\"bg-gray-800 p-6 rounded-lg shadow-lg mb-8\"><h2 class=\"text-2xl font-bold mb-4\">Generate Typing Exercise</h2><form hx-post=\"/generate-text\" hx-target=\"#typing-area\" class=\"space-y-4\"><div><label for=\"prompt\" class=\"block text-sm font-medium mb-1\">What would you like to type?</label> <input type=\"text\" id=\"prompt\" name=\"prompt\" class=\"w-full p-2 bg-gray-700 border border-gray-600 rounded focus:outline-none focus:ring-2 focus:ring-yellow-400\" placeholder=\"e.g., a Python function, a poem about coding, etc.\"></div><button type=\"submit\" class=\"w-full py-2 px-4 bg-yellow-500 hover:bg-yellow-600 text-gray-900 font-bold rounded transition\">Generate Text</button></form></div><div id=\"typing-area\" class=\"bg-gray-800 p-6 rounded-lg shadow-lg\"><p class=\"text-gray-400 text-center\">Generate a text to start typing...</p></div><div id=\"metrics\" class=\"mt-8 grid grid-cols-3 gap-4 text-center\"><div class=\"bg-gray-800 p-4 rounded-lg\"><h3 class=\"text-sm text-gray-400\">WPM</h3><p class=\"text-2xl font-bold text-yellow-400\" id=\"wpm\">0</p></div><div class=\"bg-gray-800 p-4 rounded-lg\"><h3 class=\"text-sm text-gray-400\">Accuracy</h3><p class=\"text-2xl font-bold text-yellow-400\" id=\"accuracy\">0%</p></div><div class=\"bg-gray-800 p-4 rounded-lg\"><h3 class=\"text-sm text-gray-400\">Errors</h3><p class=\"text-2xl font-bold text-yellow-400\" id=\"errors\">0</p></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
