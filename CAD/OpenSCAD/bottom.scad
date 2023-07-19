include <BOSL2/std.scad>
include <asanoha.scad>

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

asanoha_pattern_size = [
    back_asanoha_pattern_repeat.x * stand_asanoha_pattern_unit_length,
    back_asanoha_pattern_repeat.y *stand_asanoha_pattern_unit_length
];
screen_x_size = screen_count == 1 ? active_screen_area.x : screen_pcb_size.x;
module asanoha_area(asanoha_anchor = [])
{
    translate(asanoha_anchor) cuboid([ asanoha_pattern_size.x, asanoha_pattern_size.y, panel_thickness ],
                                     rounding = asanoha_area_rounding_radius,
                                     edges = [ FRONT + RIGHT, FRONT + LEFT, BACK + RIGHT, BACK + LEFT ],
                                     $fn = resolution, anchor = LEFT + BOTTOM + FRONT);
}

module panel()
{
    asanoha_anchor =
        [ panel_size.x / 2 - asanoha_pattern_size.x / 2, panel_size.y / 2 - asanoha_pattern_size.y / 2, 0 ];
    difference()
    {
        cuboid([ panel_size.x, panel_size.y, panel_thickness ], rounding = rounding_radius,
               edges = [ FRONT + RIGHT, FRONT + LEFT, BACK + RIGHT, BACK + LEFT ], $fn = resolution,
               anchor = LEFT + BOTTOM + FRONT);
        asanoha_area(asanoha_anchor);
    }
    intersection()
    {
        translate(asanoha_anchor)
            drawAsanoha(repeatXarg = back_asanoha_pattern_repeat.x, repeatYarg = back_asanoha_pattern_repeat.y);
        asanoha_area(asanoha_anchor);
    }
}

module nut_hole(thread_depth = 10, nut_depth = 1.5, nut_diameter = 2, thread_diameter = 1.5)
{
    translate([ 0, 0, thread_depth / 2 + nut_depth ])
        cylinder(h = thread_depth, r = thread_diameter / 2, center = true, $fn = resolution);
    translate([ 0, 0, nut_depth / 2 ]) cylinder(h = nut_depth, r = nut_diameter / 2, center = true, $fn = 6);
}

module screen_screw_set()
{
    translate([ 0, 0 ]) nut_hole(thread_depth = panel_thickness - m2_screw_nut_depth, nut_depth = m2_screw_nut_depth,
                                 nut_diameter = m2_screw_nut_diameter, thread_diameter = m2_screw_diameter);
    translate([ screen_pcb_screw_holes_spacing.x + screen_x_size * (screen_count - 1), 0, 0 ])
        nut_hole(thread_depth = panel_thickness - m2_screw_nut_depth, nut_depth = m2_screw_nut_depth,
                 nut_diameter = m2_screw_nut_diameter, thread_diameter = m2_screw_diameter);
    translate([ 0, screen_pcb_screw_holes_spacing.y, 0 ])
        nut_hole(thread_depth = panel_thickness - m2_screw_nut_depth, nut_depth = m2_screw_nut_depth,
                 nut_diameter = m2_screw_nut_diameter, thread_diameter = m2_screw_diameter);
    translate(
        [ screen_pcb_screw_holes_spacing.x + screen_x_size * (screen_count - 1), screen_pcb_screw_holes_spacing.y, 0 ])
        nut_hole(thread_depth = panel_thickness - m2_screw_nut_depth, nut_depth = m2_screw_nut_depth,
                 nut_diameter = m2_screw_nut_diameter, thread_diameter = m2_screw_diameter);
}

module raspberry_pi_pico_screw_set()
{
    translate([
        pi_board_screw_holes_center_from_board_edge_distance.x, pi_board_screw_holes_center_from_board_edge_distance.y
    ]) nut_hole(thread_depth = panel_thickness - m2_screw_nut_depth, nut_depth = m2_screw_nut_depth,
                nut_diameter = m2_screw_nut_diameter, thread_diameter = m2_screw_diameter);
    translate([
        pi_board_screw_holes_center_from_board_edge_distance.x,
        pi_board_size.y - pi_board_screw_holes_center_from_board_edge_distance.y
    ]) nut_hole(thread_depth = panel_thickness - m2_screw_nut_depth, nut_depth = m2_screw_nut_depth,
                nut_diameter = m2_screw_nut_diameter, thread_diameter = m2_screw_diameter);
    translate([
        pi_board_size.x - pi_board_screw_holes_center_from_board_edge_distance.x,
        pi_board_screw_holes_center_from_board_edge_distance.y
    ]) nut_hole(thread_depth = panel_thickness - m2_screw_nut_depth, nut_depth = m2_screw_nut_depth,
                nut_diameter = m2_screw_nut_diameter, thread_diameter = m2_screw_diameter);
    translate([
        pi_board_size.x - pi_board_screw_holes_center_from_board_edge_distance.x,
        pi_board_size.y - pi_board_screw_holes_center_from_board_edge_distance.y
    ]) nut_hole(thread_depth = panel_thickness - m2_screw_nut_depth, nut_depth = m2_screw_nut_depth,
                nut_diameter = m2_screw_nut_diameter, thread_diameter = m2_screw_diameter);
}

module raspberry_pi_pico()
{
    rpip_anchor = [ pi_board_offset.x, pi_board_offset.y + pcb_anchor_coords.y, 0 ];
    // translate(rpip_anchor) translate ([0,0,panel_thickness-pi_board_inset_depth])
    // cube([pi_board_size.x,pi_board_size.y,pi_board_inset_depth]);
    translate(rpip_anchor) raspberry_pi_pico_screw_set();
}

module screens()
{
    screw_dist = screen_pcb_screw_holes_center_from_board_edge_distance.y + screen_pcb_screw_holes_spacing.y -
                 active_screen_area_offset_from_corner.y - active_screen_area.y;
    screw_anchor_coords = [
        pcb_anchor_coords.x + screen_pcb_screw_holes_center_from_board_edge_distance.x,
        pcb_anchor_coords.y - screw_dist, 0
    ];
    translate(screw_anchor_coords) screen_screw_set();
}

module stand_leg(invert = false)
{
    right_chirality = invert ? LEFT : RIGHT;
    left_chirality = invert ? RIGHT : LEFT;
    part_a_edges = [ FRONT + left_chirality, BACK + left_chirality ];
    part_b_edges = [ FRONT + right_chirality, BACK + right_chirality ];
    part_a_x = invert ? rounding_radius : 0;
    part_b_x = invert ? 0 : rounding_radius;
    hinge_ball_x = invert ? -stand_hinge_ball_offset : stand_size.x + stand_hinge_ball_offset;
    translate([ 0, panel_size.y - stand_size.y, 0 ]) difference()
    {
        union()
        {
            translate([ part_a_x, 0, -stand_size.z ])
                cuboid([ rounding_radius, stand_size.y, stand_size.z ], rounding = rounding_radius,
                       edges = part_a_edges, $fn = resolution, anchor = LEFT + BOTTOM + FRONT);
            translate([ part_b_x, 0, -stand_size.z ])
                cuboid([ stand_size.x - rounding_radius, stand_size.y, stand_size.z ], rounding = stand_rounding_radius,
                       edges = part_b_edges, $fn = resolution, anchor = LEFT + BOTTOM + FRONT);
        }

        translate([ hinge_ball_x, stand_size.y / 2, -stand_size.z / 2 ])
            sphere(d = stand_hinge_ball_diameter, $fn = resolution);
    }
}

module stand()
{
    stand_leg(invert = false);
    translate([ panel_size.x - stand_size.x, 0, 0 ]) stand_leg(invert = true);
}

screen_anchor_coords =
    [ panel_size.x / 2 - screen_x_size / 2 * screen_count, panel_size.y / 2 - active_screen_area.y / 2, 0 ];

pcb_anchor_coords = [ screen_anchor_coords.x - active_screen_area_offset_from_corner.x, screen_anchor_coords.y, 0 ];

union()
{
    difference()
    {
        panel();
        screens();
        raspberry_pi_pico();
    }
    stand();
}