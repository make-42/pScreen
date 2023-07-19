include <BOSL2/std.scad>

// Start config
// Panel Size
panel_size = [ 106, 39 ]; // mm
// Pi Board
pi_board_size = [ 51, 21 ];                                        // mm
pi_board_screw_holes_center_from_board_edge_distance = [ 2, 4.8 ]; // mm
pi_board_screw_holes_diameter = 2.1;                               // mm
pi_board_offset = [ 1.75, -1.0 ];
pi_board_usb_size = [ 8.5, 4.5 ];
pi_board_usb_rounding_radius = 0.5; // mm
// pi_board_inset_depth = 0.2; //mm (keep at zero for final)

// Screen
screen_pcb_size = [ 100.5, 33.5, 13.3 ];                                // mm (+- 0.2) 13.3mm are pin headers.
screen_pcb_screw_holes_center_from_board_edge_distance = [ 2.25, 2.5 ]; // mm
screen_pcb_screw_holes_diameter = 2.8;                                  // mm
screen_pcb_screw_holes_spacing = [ 95, 28.5 ];                          // mm
active_screen_area = [ 76.78, 19.18 ];                                  // mm (old: 13.86)
active_screen_area_offset_from_corner = [ 11.36, 5.1 ];                 // mm
screen_rounding_radius = 0.5;                                           // mm
screen_count = 1;

// Screws
m2_screw_head_depth = 1;      // mm it's 2mm normally
m2_screw_head_diameter = 4.0; // mm (adding .6mm for tolerances)
m2_screw_nut_depth = 1;       // mm it's 1.6mm normally
m2_screw_nut_diameter = 5.0;  // mm (adding .6mm for tolerances)
m2_screw_diameter = 2.4;      // mm

// Asanoha
back_asanoha_pattern_repeat = [ 9, 3 ]; // mm
asanoha_area_rounding_radius = 3;       // mm

// Stand
stand_size = [ 3, 6, 6 ];
stand_rounding_radius = 1.5;                // mm
stand_hinge_ball_diameter = 2;              // mm
stand_hinge_ball_offset = 0;                // mm
stand_panel_tolerance = 0.25;               // mm
stand_panel_thickness = 1.5;                // mm
stand_panel_rounding_radius = 1.5;          // mm
stand_asanoha_pattern_unit_length = 10.175; // mm
stand_asanoha_repeat = [ 9, 3 ];

// Misc
walls_depth = 13.5;    // mm
walls_thickness = 1.5; // mm
panel_thickness = 1.5; // mm
rounding_radius = 1.5; // mm

resolution = 80;
// End config

module panel()
{
    cuboid([ panel_size.x, panel_size.y, panel_thickness ], rounding = rounding_radius,
           edges = [ FRONT + RIGHT, FRONT + LEFT, BACK + RIGHT, BACK + LEFT ], $fn = resolution,
           anchor = LEFT + BOTTOM + FRONT);
}

module screw_hole(thread_depth = 10, head_depth = 1.5, head_diameter = 2, thread_diameter = 1.5)
{
    translate([ 0, 0, panel_thickness - thread_depth / 2 - head_depth ])
        cylinder(h = thread_depth, r = thread_diameter / 2, center = true, $fn = resolution);
    translate([ 0, 0, panel_thickness - head_depth / 2 ])
        cylinder(h = head_depth, r = head_diameter / 2, center = true, $fn = resolution);
}

module screen_screw_set()
{
    translate([ 0, 0 ])
        screw_hole(thread_depth = panel_thickness - m2_screw_head_depth, head_depth = m2_screw_head_depth,
                   head_diameter = m2_screw_head_diameter, thread_diameter = m2_screw_diameter);
    translate([ screen_pcb_screw_holes_spacing.x, 0, 0 ])
        screw_hole(thread_depth = panel_thickness - m2_screw_head_depth, head_depth = m2_screw_head_depth,
                   head_diameter = m2_screw_head_diameter, thread_diameter = m2_screw_diameter);
    translate([ 0, screen_pcb_screw_holes_spacing.y, 0 ])
        screw_hole(thread_depth = panel_thickness - m2_screw_head_depth, head_depth = m2_screw_head_depth,
                   head_diameter = m2_screw_head_diameter, thread_diameter = m2_screw_diameter);
    translate([ screen_pcb_screw_holes_spacing.x, screen_pcb_screw_holes_spacing.y, 0 ])
        screw_hole(thread_depth = panel_thickness - m2_screw_head_depth, head_depth = m2_screw_head_depth,
                   head_diameter = m2_screw_head_diameter, thread_diameter = m2_screw_diameter);
}

module screens()
{
    screen_x_size = screen_count == 1 ? active_screen_area.x : screen_pcb_size.x;
    screen_anchor_coords =
        [ panel_size.x / 2 - screen_x_size / 2 * screen_count, panel_size.y / 2 - active_screen_area.y / 2, 0 ];
    pcb_anchor_coords = [ screen_anchor_coords.x - active_screen_area_offset_from_corner.x, screen_anchor_coords.y, 0 ];
    screw_dist = screen_pcb_screw_holes_center_from_board_edge_distance.y + screen_pcb_screw_holes_spacing.y -
                 active_screen_area_offset_from_corner.y - active_screen_area.y;
    screw_anchor_coords = [
        pcb_anchor_coords.x + screen_pcb_screw_holes_center_from_board_edge_distance.x,
        pcb_anchor_coords.y - screw_dist, 0
    ];
    translate(screen_anchor_coords)
        cuboid([ screen_x_size * screen_count, active_screen_area.y, panel_thickness ],
               rounding = screen_rounding_radius, edges = [ FRONT + RIGHT, FRONT + LEFT, BACK + RIGHT, BACK + LEFT ],
               $fn = resolution, anchor = LEFT + BOTTOM + FRONT);
    for (i = [0:screen_count - 1])
    {
        translate(screw_anchor_coords) translate([ screen_x_size * i, 0, 0 ]) screen_screw_set();
    }
}

difference()
{
    panel();
    screens();
}