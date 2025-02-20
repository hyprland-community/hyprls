package parser

import "image/color"

type Configuration struct {
	CustomVariables map[string]string

	General      ConfigurationGeneral      `json:"general"`
	Decoration   ConfigurationDecoration   `json:"decoration"`
	Animations   ConfigurationAnimations   `json:"animations"`
	Input        ConfigurationInput        `json:"input"`
	Gestures     ConfigurationGestures     `json:"gestures"`
	Group        ConfigurationGroup        `json:"group"`
	Misc         ConfigurationMisc         `json:"misc"`
	Binds        ConfigurationBinds        `json:"binds"`
	XWayland     ConfigurationXWayland     `json:"xwayland"`
	OpenGL       ConfigurationOpenGL       `json:"opengl"`
	Render       ConfigurationRender       `json:"render"`
	Cursor       ConfigurationCursor       `json:"cursor"`
	Ecosystem    ConfigurationEcosystem    `json:"ecosystem"`
	Experimental ConfigurationExperimental `json:"experimental"`
	Debug        ConfigurationDebug        `json:"debug"`
	Master       ConfigurationMaster       `json:"master"`
	Dwindle      ConfigurationDwindle      `json:"dwindle"`
}

type ConfigurationGeneral struct {
	// size of the border around windows
	BorderSize int `json:"border_size"`

	// disable borders for floating windows
	NoBorderOnFloating bool `json:"no_border_on_floating"`

	// gaps between windows, also supports css style gaps (top, right, bottom, left -> 5,10,15,20)
	GapsIn int `json:"gaps_in"`

	// gaps between windows and monitor edges, also supports css style gaps (top, right, bottom, left -> 5,10,15,20)
	GapsOut int `json:"gaps_out"`

	// gaps between workspaces. Stacks with gaps_out.
	GapsWorkspaces int `json:"gaps_workspaces"`

	// border color for inactive windows
	ColInactiveBorder GradientValue `json:"col.inactive_border"`

	// border color for the active window
	ColActiveBorder GradientValue `json:"col.active_border"`

	// inactive border color for window that cannot be added to a group (see denywindowfromgroup dispatcher)
	ColNogroupBorder GradientValue `json:"col.nogroup_border"`

	// active border color for window that cannot be added to a group
	ColNogroupBorderActive GradientValue `json:"col.nogroup_border_active"`

	// which layout to use. [dwindle/master]
	Layout string `json:"layout"`

	// if true, will not fall back to the next available window when moving focus in a direction where no window was found
	NoFocusFallback bool `json:"no_focus_fallback"`

	// enables resizing windows by clicking and dragging on borders and gaps
	ResizeOnBorder bool `json:"resize_on_border"`

	// extends the area around the border where you can click and drag on, only used when general:resize_on_border is on.
	ExtendBorderGrabArea int `json:"extend_border_grab_area"`

	// show a cursor icon when hovering over borders, only used when general:resize_on_border is on.
	HoverIconOnBorder bool `json:"hover_icon_on_border"`

	// master switch for allowing tearing to occur. See the Tearing page.
	AllowTearing bool `json:"allow_tearing"`

	// force floating windows to use a specific corner when being resized (1-4 going clockwise from top left, 0 to disable)
	ResizeCorner int `json:"resize_corner"`

	// Whether this configuration was autogenerated
	Autogenerated bool `json:"autogenerated"`

	Snap ConfigurationGeneralSnap `json:"snap"`
}

type ConfigurationGeneralSnap struct {
	// enable snapping for floating windows
	Enabled bool `json:"enabled"`

	// minimum gap in pixels between windows before snapping
	WindowGap int `json:"window_gap"`

	// minimum gap in pixels between window and monitor edges before snapping
	MonitorGap int `json:"monitor_gap"`

	// if true, windows snap such that only one border's worth of space is between them
	BorderOverlap bool `json:"border_overlap"`
}

type ConfigurationDecoration struct {
	// rounded corners' radius (in layout px)
	Rounding int `json:"rounding"`

	// adjusts the curve used for rounding corners, larger is smoother, 2.0 is a circle, 4.0 is a squircle. [2.0 - 10.0]
	RoundingPower float32 `json:"rounding_power"`

	// opacity of active windows. [0.0 - 1.0]
	ActiveOpacity float32 `json:"active_opacity"`

	// opacity of inactive windows. [0.0 - 1.0]
	InactiveOpacity float32 `json:"inactive_opacity"`

	// opacity of fullscreen windows. [0.0 - 1.0]
	FullscreenOpacity float32 `json:"fullscreen_opacity"`

	// enables dimming of inactive windows
	DimInactive bool `json:"dim_inactive"`

	// how much inactive windows should be dimmed [0.0 - 1.0]
	DimStrength float32 `json:"dim_strength"`

	// how much to dim the rest of the screen by when a special workspace is open. [0.0 - 1.0]
	DimSpecial float32 `json:"dim_special"`

	// how much the dimaround window rule should dim by. [0.0 - 1.0]
	DimAround float32 `json:"dim_around"`

	// a path to a custom shader to be applied at the end of rendering. See examples/screenShader.frag for an example.
	ScreenShader string `json:"screen_shader"`

	Blur   ConfigurationDecorationBlur       `json:"blur"`
	Shadow ConfigurationDecorationBlurShadow `json:"shadow"`
}

type ConfigurationDecorationBlur struct {
	// enable kawase window background blur
	Enabled bool `json:"enabled"`

	// blur size (distance)
	Size int `json:"size"`

	// the amount of passes to perform
	Passes int `json:"passes"`

	// make the blur layer ignore the opacity of the window
	IgnoreOpacity bool `json:"ignore_opacity"`

	// whether to enable further optimizations to the blur. Recommended to leave on, as it will massively improve performance.
	NewOptimizations bool `json:"new_optimizations"`

	// if enabled, floating windows will ignore tiled windows in their blur. Only available if new_optimizations is true. Will reduce overhead on floating blur significantly.
	Xray bool `json:"xray"`

	// how much noise to apply. [0.0 - 1.0]
	Noise float32 `json:"noise"`

	// contrast modulation for blur. [0.0 - 2.0]
	Contrast float32 `json:"contrast"`

	// brightness modulation for blur. [0.0 - 2.0]
	Brightness float32 `json:"brightness"`

	// Increase saturation of blurred colors. [0.0 - 1.0]
	Vibrancy float32 `json:"vibrancy"`

	// How strong the effect of vibrancy is on dark areas . [0.0 - 1.0]
	VibrancyDarkness float32 `json:"vibrancy_darkness"`

	// whether to blur behind the special workspace (note: expensive)
	Special bool `json:"special"`

	// whether to blur popups (e.g. right-click menus)
	Popups bool `json:"popups"`

	// works like ignorealpha in layer rules. If pixel opacity is below set value, will not blur. [0.0 - 1.0]
	PopupsIgnorealpha float32 `json:"popups_ignorealpha"`

	// whether to blur input methods (e.g. fcitx5)
	InputMethods bool `json:"input_methods"`

	// works like ignorealpha in layer rules. If pixel opacity is below set value, will not blur. [0.0 - 1.0]
	InputMethodsIgnorealpha float32 `json:"input_methods_ignorealpha"`
}

type ConfigurationDecorationBlurShadow struct {
	// enable drop shadows on windows
	Enabled bool `json:"enabled"`

	// Shadow range ("size") in layout px
	Range int `json:"range"`

	// in what power to render the falloff (more power, the faster the falloff) [1 - 4]
	RenderPower int `json:"render_power"`

	// if enabled, will make the shadows sharp, akin to an infinite render power
	Sharp bool `json:"sharp"`

	// if true, the shadow will not be rendered behind the window itself, only around it.
	IgnoreWindow bool `json:"ignore_window"`

	// shadow's color. Alpha dictates shadow's opacity.
	Color color.RGBA `json:"color"`

	// inactive shadow color. (if not set, will fall back to color)
	ColorInactive color.RGBA `json:"color_inactive"`

	// shadow's rendering offset.
	Offset [2]float32 `json:"offset"`

	// shadow's scale. [0.0 - 1.0]
	Scale float32 `json:"scale"`
}

type ConfigurationAnimations struct {
	// enable animations
	Enabled bool `json:"enabled"`

	// enable first launch animation
	FirstLaunchAnimation bool `json:"first_launch_animation"`
}

type ConfigurationInput struct {
	// Appropriate XKB keymap parameter. See the note below.
	KbModel string `json:"kb_model"`

	// Appropriate XKB keymap parameter
	KbLayout string `json:"kb_layout"`

	// Appropriate XKB keymap parameter
	KbVariant string `json:"kb_variant"`

	// Appropriate XKB keymap parameter
	KbOptions string `json:"kb_options"`

	// Appropriate XKB keymap parameter
	KbRules string `json:"kb_rules"`

	// If you prefer, you can use a path to your custom .xkb file.
	KbFile string `json:"kb_file"`

	// Engage numlock by default.
	NumlockByDefault bool `json:"numlock_by_default"`

	// Determines how keybinds act when multiple layouts are used. If false, keybinds will always act as if the first specified layout is active. If true, keybinds specified by symbols are activated when you type the respective symbol with the current layout.
	ResolveBindsBySym bool `json:"resolve_binds_by_sym"`

	// The repeat rate for held-down keys, in repeats per second.
	RepeatRate int `json:"repeat_rate"`

	// Delay before a held-down key is repeated, in milliseconds.
	RepeatDelay int `json:"repeat_delay"`

	// Sets the mouse input sensitivity. Value is clamped to the range -1.0 to 1.0. libinput#pointer-acceleration
	Sensitivity float32 `json:"sensitivity"`

	// Sets the cursor acceleration profile. Can be one of adaptive, flat. Can also be custom, see below. Leave empty to use libinput's default mode for your input device. libinput#pointer-acceleration [adaptive/flat/custom]
	AccelProfile string `json:"accel_profile"`

	// Force no cursor acceleration. This bypasses most of your pointer settings to get as raw of a signal as possible. Enabling this is not recommended due to potential cursor desynchronization.
	ForceNoAccel bool `json:"force_no_accel"`

	// Switches RMB and LMB
	LeftHanded bool `json:"left_handed"`

	// Sets the scroll acceleration profile, when accel_profile is set to custom. Has to be in the form <step> <points>. Leave empty to have a flat scroll curve.
	ScrollPoints string `json:"scroll_points"`

	// Sets the scroll method. Can be one of 2fg (2 fingers), edge, on_button_down, no_scroll. libinput#scrolling [2fg/edge/on_button_down/no_scroll]
	ScrollMethod string `json:"scroll_method"`

	// Sets the scroll button. Has to be an int, cannot be a string. Check wev if you have any doubts regarding the ID. 0 means default.
	ScrollButton int `json:"scroll_button"`

	// If the scroll button lock is enabled, the button does not need to be held down. Pressing and releasing the button toggles the button lock, which logically holds the button down or releases it. While the button is logically held down, motion events are converted to scroll events.
	ScrollButtonLock bool `json:"scroll_button_lock"`

	// Multiplier added to scroll movement for external mice. Note that there is a separate setting for touchpad scroll_factor.
	ScrollFactor float32 `json:"scroll_factor"`

	// Inverts scrolling direction. When enabled, scrolling moves content directly, rather than manipulating a scrollbar.
	NaturalScroll bool `json:"natural_scroll"`

	// Specify if and how cursor movement should affect window focus. See the note below. [0/1/2/3]
	FollowMouse int `json:"follow_mouse"`

	// Controls the window focus behavior when a window is closed. When set to 0, focus will shift to the next window candidate. When set to 1, focus will shift to the window under the cursor. [0/1]
	FocusOnClose int `json:"focus_on_close"`

	// If disabled, mouse focus won't switch to the hovered window unless the mouse crosses a window boundary when follow_mouse=1.
	MouseRefocus bool `json:"mouse_refocus"`

	// If enabled (1 or 2), focus will change to the window under the cursor when changing from tiled-to-floating and vice versa. If 2, focus will also follow mouse on float-to-float switches.
	FloatSwitchOverrideFocus int `json:"float_switch_override_focus"`

	// if enabled, having only floating windows in the special workspace will not block focusing windows in the regular workspace.
	SpecialFallthrough bool `json:"special_fallthrough"`

	// Handles axis events around (gaps/border for tiled, dragarea/border for floated) a focused window. 0 ignores axis events 1 sends out-of-bound coordinates 2 fakes pointer coordinates to the closest point inside the window 3 warps the cursor to the closest point inside the window
	OffWindowAxisEvents int `json:"off_window_axis_events"`

	// Emulates discrete scrolling from high resolution scrolling events. 0 disables it, 1 enables handling of non-standard events only, and 2 force enables all scroll wheel events to be handled
	EmulateDiscreteScroll int `json:"emulate_discrete_scroll"`
}

type ConfigurationGestures struct {
	// enable workspace swipe gesture on touchpad
	WorkspaceSwipe bool `json:"workspace_swipe"`

	// how many fingers for the touchpad gesture
	WorkspaceSwipeFingers int `json:"workspace_swipe_fingers"`

	// if enabled, workspace_swipe_fingers is considered the minimum number of fingers to swipe
	WorkspaceSwipeMinFingers bool `json:"workspace_swipe_min_fingers"`

	// in px, the distance of the touchpad gesture
	WorkspaceSwipeDistance int `json:"workspace_swipe_distance"`

	// enable workspace swiping from the edge of a touchscreen
	WorkspaceSwipeTouch bool `json:"workspace_swipe_touch"`

	// invert the direction (touchpad only)
	WorkspaceSwipeInvert bool `json:"workspace_swipe_invert"`

	// invert the direction (touchscreen only)
	WorkspaceSwipeTouchInvert bool `json:"workspace_swipe_touch_invert"`

	// minimum speed in px per timepoint to force the change ignoring cancel_ratio. Setting to 0 will disable this mechanic.
	WorkspaceSwipeMinSpeedToForce int `json:"workspace_swipe_min_speed_to_force"`

	// how much the swipe has to proceed in order to commence it. (0.7 -> if > 0.7 * distance, switch, if less, revert) [0.0 - 1.0]
	WorkspaceSwipeCancelRatio float32 `json:"workspace_swipe_cancel_ratio"`

	// whether a swipe right on the last workspace should create a new one.
	WorkspaceSwipeCreateNew bool `json:"workspace_swipe_create_new"`

	// if enabled, switching direction will be locked when you swipe past the direction_lock_threshold (touchpad only).
	WorkspaceSwipeDirectionLock bool `json:"workspace_swipe_direction_lock"`

	// in px, the distance to swipe before direction lock activates (touchpad only).
	WorkspaceSwipeDirectionLockThreshold int `json:"workspace_swipe_direction_lock_threshold"`

	// if enabled, swiping will not clamp at the neighboring workspaces but continue to the further ones.
	WorkspaceSwipeForever bool `json:"workspace_swipe_forever"`

	// if enabled, swiping will use the r prefix instead of the m prefix for finding workspaces.
	WorkspaceSwipeUseR bool `json:"workspace_swipe_use_r"`
}

type ConfigurationGroup struct {
	// whether new windows will be automatically grouped into the focused unlocked group. Note: if you want to disable auto_group only for specific windows, use the "group barred" window rule instead.
	AutoGroup bool `json:"auto_group"`

	// whether new windows in a group spawn after current or at group tail
	InsertAfterCurrent bool `json:"insert_after_current"`

	// whether Hyprland should focus on the window that has just been moved out of the group
	FocusRemovedWindow bool `json:"focus_removed_window"`

	// whether dragging a window into a unlocked group will merge them. Options: 0 (disabled), 1 (enabled), 2 (only when dragging into the groupbar)
	DragIntoGroup int `json:"drag_into_group"`

	// whether window groups can be dragged into other groups
	MergeGroupsOnDrag bool `json:"merge_groups_on_drag"`

	// whether one group will be merged with another when dragged into its groupbar
	MergeGroupsOnGroupbar bool `json:"merge_groups_on_groupbar"`

	// whether dragging a floating window into a tiled window groupbar will merge them
	MergeFloatedIntoTiledOnGroupbar bool `json:"merge_floated_into_tiled_on_groupbar"`

	// whether using movetoworkspace[silent] will merge the window into the workspace's solitary unlocked group
	GroupOnMovetoworkspace bool `json:"group_on_movetoworkspace"`

	// active group border color
	ColBorderActive GradientValue `json:"col.border_active"`

	// inactive (out of focus) group border color
	ColBorderInactive GradientValue `json:"col.border_inactive"`

	// active locked group border color
	ColBorderLockedActive GradientValue `json:"col.border_locked_active"`

	// inactive locked group border color
	ColBorderLockedInactive GradientValue `json:"col.border_locked_inactive"`

	Groupbar ConfigurationGroupGroupbar `json:"groupbar"`
}

type ConfigurationGroupGroupbar struct {
	// enables groupbars
	Enabled bool `json:"enabled"`

	// font used to display groupbar titles, use misc:font_family if not specified
	FontFamily string `json:"font_family"`

	// font size of groupbar title
	FontSize int `json:"font_size"`

	// enables gradients
	Gradients bool `json:"gradients"`

	// height of the groupbar
	Height int `json:"height"`

	// render the groupbar as a vertical stack
	Stacked bool `json:"stacked"`

	// sets the decoration priority for groupbars
	Priority int `json:"priority"`

	// whether to render titles in the group bar decoration
	RenderTitles bool `json:"render_titles"`

	// whether scrolling in the groupbar changes group active window
	Scrolling bool `json:"scrolling"`

	// controls the group bar text color
	TextColor color.RGBA `json:"text_color"`

	// active group bar background color
	ColActive GradientValue `json:"col.active"`

	// inactive (out of focus) group bar background color
	ColInactive GradientValue `json:"col.inactive"`

	// active locked group bar background color
	ColLockedActive GradientValue `json:"col.locked_active"`

	// inactive locked group bar background color
	ColLockedInactive GradientValue `json:"col.locked_inactive"`
}

type ConfigurationMisc struct {
	// disables the random Hyprland logo / anime girl background. :(
	DisableHyprlandLogo bool `json:"disable_hyprland_logo"`

	// disables the Hyprland splash rendering. (requires a monitor reload to take effect)
	DisableSplashRendering bool `json:"disable_splash_rendering"`

	// Changes the color of the splash text (requires a monitor reload to take effect).
	ColSplash color.RGBA `json:"col.splash"`

	// Set the global default font to render the text including debug fps/notification, config error messages and etc., selected from system fonts.
	FontFamily string `json:"font_family"`

	// Changes the font used to render the splash text, selected from system fonts (requires a monitor reload to take effect).
	SplashFontFamily string `json:"splash_font_family"`

	// Enforce any of the 3 default wallpapers. Setting this to 0 or 1 disables the anime background. -1 means "random". [-1/0/1/2]
	ForceDefaultWallpaper int `json:"force_default_wallpaper"`

	// controls the VFR status of Hyprland. Heavily recommended to leave enabled to conserve resources.
	Vfr bool `json:"vfr"`

	// controls the VRR (Adaptive Sync) of your monitors. 0 - off, 1 - on, 2 - fullscreen only [0/1/2]
	Vrr int `json:"vrr"`

	// If DPMS is set to off, wake up the monitors if the mouse moves.
	MouseMoveEnablesDpms bool `json:"mouse_move_enables_dpms"`

	// If DPMS is set to off, wake up the monitors if a key is pressed.
	KeyPressEnablesDpms bool `json:"key_press_enables_dpms"`

	// Will make mouse focus follow the mouse when drag and dropping. Recommended to leave it enabled, especially for people using focus follows mouse at 0.
	AlwaysFollowOnDnd bool `json:"always_follow_on_dnd"`

	// If true, will make keyboard-interactive layers keep their focus on mouse move (e.g. wofi, bemenu)
	LayersHogKeyboardFocus bool `json:"layers_hog_keyboard_focus"`

	// If true, will animate manual window resizes/moves
	AnimateManualResizes bool `json:"animate_manual_resizes"`

	// If true, will animate windows being dragged by mouse, note that this can cause weird behavior on some curves
	AnimateMouseWindowdragging bool `json:"animate_mouse_windowdragging"`

	// If true, the config will not reload automatically on save, and instead needs to be reloaded with hyprctl reload. Might save on battery.
	DisableAutoreload bool `json:"disable_autoreload"`

	// Enable window swallowing
	EnableSwallow bool `json:"enable_swallow"`

	// The class regex to be used for windows that should be swallowed (usually, a terminal). To know more about the list of regex which can be used use this cheatsheet.
	SwallowRegex string `json:"swallow_regex"`

	// The title regex to be used for windows that should not be swallowed by the windows specified in swallow_regex  (e.g. wev). The regex is matched against the parent (e.g. Kitty) window's title on the assumption that it changes to whatever process it's running.
	SwallowExceptionRegex string `json:"swallow_exception_regex"`

	// Whether Hyprland should focus an app that requests to be focused (an activate request)
	FocusOnActivate bool `json:"focus_on_activate"`

	// Whether mouse moving into a different monitor should focus it
	MouseMoveFocusesMonitor bool `json:"mouse_move_focuses_monitor"`

	// [Warning: buggy] starts rendering before your monitor displays a frame in order to lower latency
	RenderAheadOfTime bool `json:"render_ahead_of_time"`

	// how many ms of safezone to add to rendering ahead of time. Recommended 1-2.
	RenderAheadSafezone int `json:"render_ahead_safezone"`

	// if true, will allow you to restart a lockscreen app in case it crashes (red screen of death)
	AllowSessionLockRestore bool `json:"allow_session_lock_restore"`

	// change the background color. (requires enabled disable_hyprland_logo)
	BackgroundColor color.RGBA `json:"background_color"`

	// close the special workspace if the last window is removed
	CloseSpecialOnEmpty bool `json:"close_special_on_empty"`

	// if there is a fullscreen or maximized window, decide whether a new tiled window opened should replace it, stay behind or disable the fullscreen/maximized state. 0 - behind, 1 - takes over, 2 - unfullscreen/unmaxize [0/1/2]
	NewWindowTakesOverFullscreen int `json:"new_window_takes_over_fullscreen"`

	// if true, closing a fullscreen window makes the next focused window fullscreen
	ExitWindowRetainsFullscreen bool `json:"exit_window_retains_fullscreen"`

	// if enabled, windows will open on the workspace they were invoked on. 0 - disabled, 1 - single-shot, 2 - persistent (all children too)
	InitialWorkspaceTracking int `json:"initial_workspace_tracking"`

	// whether to enable middle-click-paste (aka primary selection)
	MiddleClickPaste bool `json:"middle_click_paste"`

	// the maximum limit for renderunfocused windows' fps in the background (see also Window-Rules - renderunfocused)
	RenderUnfocusedFps int `json:"render_unfocused_fps"`

	// disable the warning if XDG environment is externally managed
	DisableXdgEnvChecks bool `json:"disable_xdg_env_checks"`

	// disable the warning if hyprland-qtutils is not installed
	DisableHyprlandQtutilsCheck bool `json:"disable_hyprland_qtutils_check"`

	// the delay in ms after the lockdead screen appears if the lock screen did not appear after a lock event occurred
	LockdeadScreenDelay int `json:"lockdead_screen_delay"`
}

type ConfigurationBinds struct {
	// if disabled, will not pass the mouse events to apps / dragging windows around if a keybind has been triggered.
	PassMouseWhenBound bool `json:"pass_mouse_when_bound"`

	// in ms, how many ms to wait after a scroll event to allow passing another one for the binds.
	ScrollEventDelay int `json:"scroll_event_delay"`

	// If enabled, an attempt to switch to the currently focused workspace will instead switch to the previous workspace. Akin to i3's auto_back_and_forth.
	WorkspaceBackAndForth bool `json:"workspace_back_and_forth"`

	// If enabled, workspaces don't forget their previous workspace, so cycles can be created by switching to the first workspace in a sequence, then endlessly going to the previous workspace.
	AllowWorkspaceCycles bool `json:"allow_workspace_cycles"`

	// Whether switching workspaces should center the cursor on the workspace (0) or on the last active window for that workspace (1)
	WorkspaceCenterOn int `json:"workspace_center_on"`

	// sets the preferred focus finding method when using focuswindow/movewindow/etc with a direction. 0 - history (recent have priority), 1 - length (longer shared edges have priority)
	FocusPreferredMethod int `json:"focus_preferred_method"`

	// If enabled, dispatchers like moveintogroup, moveoutofgroup and movewindoworgroup will ignore lock per group.
	IgnoreGroupLock bool `json:"ignore_group_lock"`

	// If enabled, when on a fullscreen window, movefocus will cycle fullscreen, if not, it will move the focus in a direction.
	MovefocusCyclesFullscreen bool `json:"movefocus_cycles_fullscreen"`

	// If enabled, when in a grouped window, movefocus will cycle windows in the groups first, then at each ends of tabs, it'll move on to other windows/groups
	MovefocusCyclesGroupfirst bool `json:"movefocus_cycles_groupfirst"`

	// If enabled, apps that request keybinds to be disabled (e.g. VMs) will not be able to do so.
	DisableKeybindGrabbing bool `json:"disable_keybind_grabbing"`

	// If enabled, moving a window or focus over the edge of a monitor with a direction will move it to the next monitor in that direction.
	WindowDirectionMonitorFallback bool `json:"window_direction_monitor_fallback"`

	// If enabled, Allow fullscreen to pinned windows, and restore their pinned status afterwards
	AllowPinFullscreen bool `json:"allow_pin_fullscreen"`
}

type ConfigurationXWayland struct {
	// allow running applications using X11
	Enabled bool `json:"enabled"`

	// uses the nearest neighbor filtering for xwayland apps, making them pixelated rather than blurry
	UseNearestNeighbor bool `json:"use_nearest_neighbor"`

	// forces a scale of 1 on xwayland windows on scaled displays.
	ForceZeroScaling bool `json:"force_zero_scaling"`
}

type ConfigurationOpenGL struct {
	// reduces flickering on nvidia at the cost of possible frame drops on lower-end GPUs. On non-nvidia, this is ignored.
	NvidiaAntiFlicker bool `json:"nvidia_anti_flicker"`

	// forces introspection at all times. Introspection is aimed at reducing GPU usage in certain cases, but might cause graphical glitches on nvidia. 0 - nothing, 1 - force always on, 2 - force always on if nvidia
	ForceIntrospection int `json:"force_introspection"`
}

type ConfigurationRender struct {
	// Whether to enable explicit sync support. Requires a hyprland restart. 0 - no, 1 - yes, 2 - auto based on the gpu driver
	ExplicitSync int `json:"explicit_sync"`

	// Whether to enable explicit sync support for the KMS layer. Requires explicit_sync to be enabled. 0 - no, 1 - yes, 2 - auto based on the gpu driver
	ExplicitSyncKms int `json:"explicit_sync_kms"`

	// Enables direct scanout. Direct scanout attempts to reduce lag when there is only one fullscreen application on a screen (e.g. game). It is also recommended to set this to false if the fullscreen application shows graphical glitches.
	DirectScanout bool `json:"direct_scanout"`

	// Whether to expand undersized textures along the edge, or rather stretch the entire texture.
	ExpandUndersizedTextures bool `json:"expand_undersized_textures"`

	// Disables back buffer and bottom layer rendering.
	XpMode bool `json:"xp_mode"`

	// Whether to enable a fade animation for CTM changes (hyprsunset). 2 means "auto" which disables them on Nvidia.
	CtmAnimation int `json:"ctm_animation"`

	// Allow early buffer release event. Fixes stuttering and missing frames for some apps. May cause graphical glitches and memory leaks in others.
	AllowEarlyBufferRelease bool `json:"allow_early_buffer_release"`
}

type ConfigurationCursor struct {
	// sync xcursor theme with gsettings, it applies cursor-theme and cursor-size on theme load to gsettings making most CSD gtk based clients use same xcursor theme and size.
	SyncGsettingsTheme bool `json:"sync_gsettings_theme"`

	// disables hardware cursors.
	NoHardwareCursors bool `json:"no_hardware_cursors"`

	// disables scheduling new frames on cursor movement for fullscreen apps with VRR enabled to avoid framerate spikes (requires no_hardware_cursors = true)
	NoBreakFsVrr bool `json:"no_break_fs_vrr"`

	// minimum refresh rate for cursor movement when no_break_fs_vrr is active. Set to minimum supported refresh rate or higher
	MinRefreshRate int `json:"min_refresh_rate"`

	// the padding, in logical px, between screen edges and the cursor
	HotspotPadding int `json:"hotspot_padding"`

	// in seconds, after how many seconds of cursor's inactivity to hide it. Set to 0 for never.
	InactiveTimeout float32 `json:"inactive_timeout"`

	// if true, will not warp the cursor in many cases (focusing, keybinds, etc)
	NoWarps bool `json:"no_warps"`

	// When a window is refocused, the cursor returns to its last position relative to that window, rather than to the centre.
	PersistentWarps bool `json:"persistent_warps"`

	// Move the cursor to the last focused window after changing the workspace. Options: 0 (Disabled), 1 (Enabled), 2 (Force - ignores cursor:no_warps option)
	WarpOnChangeWorkspace int `json:"warp_on_change_workspace"`

	// the name of a default monitor for the cursor to be set to on startup (see hyprctl monitors for names)
	DefaultMonitor string `json:"default_monitor"`

	// the factor to zoom by around the cursor. Like a magnifying glass. Minimum 1.0 (meaning no zoom)
	ZoomFactor float32 `json:"zoom_factor"`

	// whether the zoom should follow the cursor rigidly (cursor is always centered if it can be) or loosely
	ZoomRigid bool `json:"zoom_rigid"`

	// whether to enable hyprcursor support
	EnableHyprcursor bool `json:"enable_hyprcursor"`

	// Hides the cursor when you press any key until the mouse is moved.
	HideOnKeyPress bool `json:"hide_on_key_press"`

	// Hides the cursor when the last input was a touch input until a mouse input is done.
	HideOnTouch bool `json:"hide_on_touch"`

	// Makes HW cursors use a CPU buffer. Required on Nvidia to have HW cursors. 0 - off, 1 - on, 2 - auto (nvidia only)
	UseCpuBuffer int `json:"use_cpu_buffer"`

	// Warp the cursor back to where it was after using a non-mouse input to move it, and then returning back to mouse.
	WarpBackAfterNonMouseInput bool `json:"warp_back_after_non_mouse_input"`
}

type ConfigurationEcosystem struct {
	// disable the popup that shows up when you update hyprland to a new version.
	NoUpdateNews bool `json:"no_update_news"`

	// disable the popup that shows up twice a year encouraging to donate.
	NoDonationNag bool `json:"no_donation_nag"`
}

type ConfigurationExperimental struct {
	// force wide color gamut for all supported outputs
	WideColorGamut bool `json:"wide_color_gamut"`

	// force static hdr for all supported outputs (for testing only, will result in oversaturated colors)
	Hdr bool `json:"hdr"`

	// enable color management protocol
	XxColorManagementV4 bool `json:"xx_color_management_v4"`
}

type ConfigurationDebug struct {
	// print the debug performance overlay. Disable VFR for accurate results.
	Overlay bool `json:"overlay"`

	// (epilepsy warning!) flash areas updated with damage tracking
	DamageBlink bool `json:"damage_blink"`

	// disable logging to a file
	DisableLogs bool `json:"disable_logs"`

	// disables time logging
	DisableTime bool `json:"disable_time"`

	// redraw only the needed bits of the display. Do not change. (default: full - 2) monitor - 1, none - 0
	DamageTracking int `json:"damage_tracking"`

	// enables logging to stdout
	EnableStdoutLogs bool `json:"enable_stdout_logs"`

	// set to 1 and then back to 0 to crash Hyprland.
	ManualCrash int `json:"manual_crash"`

	// if true, do not display config file parsing errors.
	SuppressErrors bool `json:"suppress_errors"`

	// sets the timeout in seconds for watchdog to abort processing of a signal of the main thread. Set to 0 to disable.
	WatchdogTimeout int `json:"watchdog_timeout"`

	// disables verification of the scale factors. Will result in pixel alignment and rounding errors.
	DisableScaleChecks bool `json:"disable_scale_checks"`

	// limits the number of displayed config file parsing errors.
	ErrorLimit int `json:"error_limit"`

	// sets the position of the error bar. top - 0, bottom - 1
	ErrorPosition int `json:"error_position"`

	// enables colors in the stdout logs.
	ColoredStdoutLogs bool `json:"colored_stdout_logs"`

	// enables render pass debugging.
	Pass bool `json:"pass"`
}

type ConfigurationMaster struct {
	// enable adding additional master windows in a horizontal split style
	AllowSmallSplit bool `json:"allow_small_split"`

	// the scale of the special workspace windows. [0.0 - 1.0]
	SpecialScaleFactor float32 `json:"special_scale_factor"`

	// the size as a percentage of the master window, for example mfact = 0.70 would mean 70% of the screen will be the master window, and 30% the slave [0.0 - 1.0]
	Mfact float32 `json:"mfact"`

	// master: new window becomes master; slave: new windows are added to slave stack; inherit: inherit from focused window
	NewStatus string `json:"new_status"`

	// whether a newly open window should be on the top of the stack
	NewOnTop bool `json:"new_on_top"`

	// before, after: place new window relative to the focused window; none: place new window according to the value of new_on_top.
	NewOnActive string `json:"new_on_active"`

	// default placement of the master area, can be left, right, top, bottom or center
	Orientation string `json:"orientation"`

	// inherit fullscreen status when cycling/swapping to another window (e.g. monocle layout)
	InheritFullscreen bool `json:"inherit_fullscreen"`

	// when using orientation=center, make the master window centered only when at least this many slave windows are open. (Set 0 to always_center_master)
	SlaveCountForCenterMaster int `json:"slave_count_for_center_master"`

	// set if the slaves should appear on right of master when slave_count_for_center_master > 2
	CenterMasterSlavesOnRight bool `json:"center_master_slaves_on_right"`

	// if enabled, resizing direction will be determined by the mouse's position on the window (nearest to which corner). Else, it is based on the window's tiling position.
	SmartResizing bool `json:"smart_resizing"`

	// when enabled, dragging and dropping windows will put them at the cursor position. Otherwise, when dropped at the stack side, they will go to the top/bottom of the stack depending on new_on_top.
	DropAtCursor bool `json:"drop_at_cursor"`
}

type ConfigurationDwindle struct {
	// enable pseudotiling. Pseudotiled windows retain their floating size when tiled.
	Pseudotile bool `json:"pseudotile"`

	// 0 -> split follows mouse, 1 -> always split to the left (new = left or top) 2 -> always split to the right (new = right or bottom)
	ForceSplit int `json:"force_split"`

	// if enabled, the split (side/top) will not change regardless of what happens to the container.
	PreserveSplit bool `json:"preserve_split"`

	// if enabled, allows a more precise control over the window split direction based on the cursor's position. The window is conceptually divided into four triangles, and cursor's triangle determines the split direction. This feature also turns on preserve_split.
	SmartSplit bool `json:"smart_split"`

	// if enabled, resizing direction will be determined by the mouse's position on the window (nearest to which corner). Else, it is based on the window's tiling position.
	SmartResizing bool `json:"smart_resizing"`

	// if enabled, makes the preselect direction persist until either this mode is turned off, another direction is specified, or a non-direction is specified (anything other than l,r,u/t,d/b)
	PermanentDirectionOverride bool `json:"permanent_direction_override"`

	// specifies the scale factor of windows on the special workspace [0 - 1]
	SpecialScaleFactor float32 `json:"special_scale_factor"`

	// specifies the auto-split width multiplier
	SplitWidthMultiplier float32 `json:"split_width_multiplier"`

	// whether to prefer the active window or the mouse position for splits
	UseActiveForSplits bool `json:"use_active_for_splits"`

	// the default split ratio on window open. 1 means even 50/50 split. [0.1 - 1.9]
	DefaultSplitRatio float32 `json:"default_split_ratio"`

	// specifies which window will receive the larger half of a split. positional - 0, current window - 1, opening window - 2 [0/1/2]
	SplitBias int `json:"split_bias"`
}
